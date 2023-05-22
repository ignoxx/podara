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

func (s *Server) saveAudio(imageFile multipart.File, filename string) (string, error) {
	out, err := os.Create(s.audioDir + "/" + filename)

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
func (s *Server) handleCreateEpisode(w http.ResponseWriter, r *http.Request) error {
	podcastId := mux.Vars(r)["podcast_id"]

	if podcastId == "" {
		return errors.New("podcast_id is missing")
	}

	// check if podcast exists
	podcast, err := s.store.GetPodcastByID(podcastId)

	if err != nil {
		return err
	}

	episode := &types.Episode{PodcastId: podcast.Id}
	r.ParseMultipartForm(32 << 20) // limit to 32MB
	imageFile, handler, err := r.FormFile("cover_image")

	if err != nil {
		return err
	}

	defer imageFile.Close()

	audioFile, handler, err := r.FormFile("audio_file")

	if err != nil {
		return err
	}

	defer audioFile.Close()

	episode.Title = r.FormValue("title")
	episode.Description = r.FormValue("description")

	sanitizedTitle := strings.Replace(episode.Title, " ", "_", -1)
	imageFileExtension := strings.Split(handler.Filename, ".")[len(strings.Split(handler.Filename, "."))-1]
	audioFileExtension := strings.Split(handler.Filename, ".")[len(strings.Split(handler.Filename, "."))-1]

	imageFileName, err := s.saveImage(imageFile, fmt.Sprintf("episode_cover_image_%s.%s", sanitizedTitle, imageFileExtension))

	if err != nil {
		return err
	}

	audioFileName, err := s.saveAudio(audioFile, fmt.Sprintf("audio_file_%s.%s", sanitizedTitle, audioFileExtension))

	if err != nil {
		return err
	}

	episode.CoverImageUrl = imageFileName
	episode.AudioUrl = audioFileName

	s.store.CreateEpisode(episode)

	return WriteJSON(w, http.StatusOK, episode)
}

func (s *Server) handleGetEpisode(w http.ResponseWriter, r *http.Request) error {
	episodeId := mux.Vars(r)["episode_id"]
	podcastId := mux.Vars(r)["podcast_id"]

	episode, err := s.store.GetEpisodeByID(podcastId, episodeId)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, episode)
}

func (s *Server) handleUpdateEpisode(w http.ResponseWriter, r *http.Request) error {
	episodeId := mux.Vars(r)["episode_id"]
	podcastId := mux.Vars(r)["podcast_id"]

	episode, err := s.store.GetEpisodeByID(podcastId, episodeId)

	if err != nil {
		return err
	}

	if err := json.NewDecoder(r.Body).Decode(episode); err != nil {
		return err
	}

	_, err = s.store.UpdateEpisode(podcastId, episode)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, episode)
}

func (s *Server) handleDeleteEpisode(w http.ResponseWriter, r *http.Request) error {
	episodeId := mux.Vars(r)["episode_id"]
	podcastId := mux.Vars(r)["podcast_id"]

	episode, err := s.store.GetEpisodeByID(podcastId, episodeId)

	if err != nil {
		return err
	}

	err = s.store.DeleteEpisode(podcastId, episode.Id)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, episode)
}
