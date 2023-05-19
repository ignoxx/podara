package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ignoxx/podara/poc3/types"
)

func (s *Server) handleCreatePodcast(w http.ResponseWriter, r *http.Request) error {
	podcast := &types.Podcast{}

	if err := json.NewDecoder(r.Body).Decode(podcast); err != nil {
		return err
	}

	userClaims, err := getJwtClaims(r)

	if err != nil {
		return err
	}

	podcast.UserId = userClaims.Id

	_, err = s.store.CreatePodcast(podcast)

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
