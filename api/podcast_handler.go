package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
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
	podcast, err := s.store.GetPodcastByID(podcastId)

	if err != nil {
		return err
	}

	// TODO: generate rss feed

	return WriteJSON(w, http.StatusOK, podcast)
}
