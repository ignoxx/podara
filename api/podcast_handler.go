package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	rss "github.com/ignoxx/podara/poc3/pkg/podcast"
	"github.com/ignoxx/podara/poc3/types"
)

func (s *Server) saveImage(imageFile multipart.File, filename string) (string, error) {
	out, err := os.Create(s.imageDir + "/" + filename)

	if err != nil {
		return "", err
	}

	defer out.Close()

	_, err = io.Copy(out, imageFile)

	if err != nil {
		return "", err
	}

	return filename, nil
}

func (s *Server) handleCreatePodcast(w http.ResponseWriter, r *http.Request) error {
	podcast := &types.Podcast{}

	r.ParseMultipartForm(32 << 20) // limit to 32MB
	imageFile, handler, err := r.FormFile("cover_image")

	if err != nil {
		return err
	}

	defer imageFile.Close()

	// validate form values
	if r.FormValue("title") == "" || r.FormValue("description") == "" {
		return errors.New("missing title and description")
	}

	podcast.Title = r.FormValue("title")
	podcast.Description = r.FormValue("description")

	sanitizedTitle := strings.Replace(podcast.Title, " ", "_", -1)
	fileExtension := strings.Split(handler.Filename, ".")[len(strings.Split(handler.Filename, "."))-1]
	filename, err := s.saveImage(imageFile, fmt.Sprintf("cover_image_%s.%s", sanitizedTitle, fileExtension))

	if err != nil {
		return err
	}

	podcast.CoverImageUrl = filename

	userClaims, err := getJwtClaims(r)

	if err != nil {
		println("jwt claims error")
		return err
	}

	podcast.UserId = userClaims.Id

	podcast, err = s.store.CreatePodcast(podcast)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, podcast)
}

func (s *Server) handleGetAllPodcasts(w http.ResponseWriter, r *http.Request) error {
	podcasts, err := s.store.GetAllPodcasts()

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, podcasts)
}
func (s *Server) handleGetPodcast(w http.ResponseWriter, r *http.Request) error {
	podcastId := mux.Vars(r)["podcast_id"]
	podcast, err := s.store.GetPodcastByID(podcastId)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, podcast)
}

func (s *Server) handleUpdatePodcast(w http.ResponseWriter, r *http.Request) error {
	podcastId := mux.Vars(r)["podcast_id"]
	podcast, err := s.store.GetPodcastByID(podcastId)

	if err != nil {
		return err
	}

	if err := json.NewDecoder(r.Body).Decode(podcast); err != nil {
		return err
	}

	_, err = s.store.UpdatePodcast(podcast)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, podcast)
}

func (s *Server) handleDeletePodcast(w http.ResponseWriter, r *http.Request) error {
	podcastId := mux.Vars(r)["podcast_id"]
	err := s.store.DeletePodcast(podcastId)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, nil)
}

func (s *Server) handleGetPodcastRss(w http.ResponseWriter, r *http.Request) error {
	podcastId := mux.Vars(r)["podcast_id"]
	podcast, episodes, err := s.store.GetPodcastAndEpisodesByPodcastID(podcastId)

	if err != nil {
		return err
	}

	domain := "http://ec2-18-157-73-46.eu-central-1.compute.amazonaws.com"

	feed := rss.New(
		podcast.Title,
		domain,
		podcast.Description,
		SqliteDatetimeToRssDatetime(&podcast.CreatedAt),
		SqliteDatetimeToRssDatetime(&podcast.UpdatedAt),
	)

	feed.Generator = "Podara v1.0"
	feed.IExplicit = "no"
	feed.AddAuthor("John Doe", "example@mail.com")
	feed.AddAtomLink(domain + "/api/v1/podcast/" + podcastId + "/rss.xml")
	feed.AddImage(domain + "/" + podcast.CoverImageUrl)
	feed.AddSummary("This is the podcast description")
	feed.AddCategory("Technology", []string{})
	feed.Copyright = "podara copyright"

	for _, episode := range episodes {
		e := rss.Item{
			Title:       episode.Title,
			Description: episode.Description,
			Link:        domain + "/" + episode.AudioUrl,
			PubDate:     SqliteDatetimeToRssDatetime(&episode.CreatedAt),
		}

		e.AddImage(domain + "/" + episode.CoverImageUrl)
		e.AddSummary("This is the episode description")
		e.AddEnclosure(domain+"/"+episode.AudioUrl, rss.MP3, 100)
		e.AddDuration(60)

		if _, err := feed.AddItem(e); err != nil {
			fmt.Println(e.Title, ": ", err)
			return err
		}
	}

	w.Header().Set("Last-Modified", SqliteDatetimeToRssDatetime(&podcast.UpdatedAt).Format(time.RFC1123))
	// add eTag header: https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/ETag
	w.Header().Set("ETag", "123456789")
	w.Header().Set("Content-Length", strconv.Itoa(len(feed.String())))

	return feed.Encode(w)
}

func SqliteDatetimeToRssDatetime(datetime *string) *time.Time {
	t, err := time.Parse("2006-01-02T15:04:05Z", *datetime)

	if err != nil {
		return nil
	}

	return &t
}
