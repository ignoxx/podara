package api

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ignoxx/podara/poc3/storage"
)

type apiFunc func(http.ResponseWriter, *http.Request) error

type ApiError struct {
	Error string
}

func makeHTTPHandleFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			WriteJSON(w, http.StatusBadRequest, ApiError{Error: err.Error()})
		}
	}
}

type Server struct {
	listenAddr string
	store      storage.Storage
}

func NewServer(listenAddr string, store storage.Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/register", makeHTTPHandleFunc(s.handleCreateUser)).Methods("POST")
	router.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods("POST")
	router.HandleFunc("/testAuth", withAuth(makeHTTPHandleFunc(s.handleTestAuth))).Methods("GET")

	router.HandleFunc("/podcast", makeHTTPHandleFunc(s.handleCreatePodcast)).Methods("POST")
	router.HandleFunc("/podcast/{id}", makeHTTPHandleFunc(s.handleGetPodcast)).Methods("GET")
	router.HandleFunc("/podcast/{id}", makeHTTPHandleFunc(s.handleUpdatePodcast)).Methods("PUT")
	router.HandleFunc("/podcast/{id}", makeHTTPHandleFunc(s.handleDeletePodcast)).Methods("DELETE")
	router.HandleFunc("/podcast/{id}/rss.xml", makeHTTPHandleFunc(s.handleGetPodcastRss)).Methods("GET")

	router.HandleFunc("/episode", makeHTTPHandleFunc(s.handleCreateEpisode)).Methods("POST")
	router.HandleFunc("/episode/{id}", makeHTTPHandleFunc(s.handleGetEpisode)).Methods("GET")
	router.HandleFunc("/episode/{id}", makeHTTPHandleFunc(s.handleUpdateEpisode)).Methods("PUT")
	router.HandleFunc("/episode/{id}", makeHTTPHandleFunc(s.handleDeleteEpisode)).Methods("DELETE")

	return http.ListenAndServe(s.listenAddr, router)
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

// create handleTestAuth
func (s *Server) handleTestAuth(w http.ResponseWriter, r *http.Request) error {
	return nil
}
