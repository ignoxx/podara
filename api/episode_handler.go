package api

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
	"github.com/ignoxx/podara/poc3/types"
)

func (s *Server) handleCreateEpisode(w http.ResponseWriter, r *http.Request) error {
	var episode *types.Episode

	if err := json.NewDecoder(r.Body).Decode(episode); err != nil {
		return err
	}

	episode.Id = uuid.NewString()

	_, err := s.store.CreateEpisode(episode)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, episode)
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
