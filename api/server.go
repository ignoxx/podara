package api

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ignoxx/podara/poc3/storage"
)

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

    v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/register", makeHTTPHandleFunc(s.handleCreateUser)).Methods("POST")
	v1.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods("POST")

	v1.HandleFunc("/podcast", withAuth(makeHTTPHandleFunc(s.handleCreatePodcast))).Methods("POST")
	v1.HandleFunc("/podcast", withAuth(makeHTTPHandleFunc(s.handleGetAllPodcasts))).Methods("GET")
	v1.HandleFunc("/podcast/{podcastID}", withAuth(makeHTTPHandleFunc(s.handleGetPodcast))).Methods("GET")
	v1.HandleFunc("/podcast/{podcastID}", withAuth(makeHTTPHandleFunc(s.handleUpdatePodcast))).Methods("PUT")
	v1.HandleFunc("/podcast/{podcastID}", withAuth(makeHTTPHandleFunc(s.handleDeletePodcast))).Methods("DELETE")

	v1.HandleFunc("/podcast/{podcastID}/rss.xml", makeHTTPHandleFunc(s.handleGetPodcastRss)).Methods("GET")

	v1.HandleFunc("/podcast/{podcastID}/episode", withAuth(makeHTTPHandleFunc(s.handleCreateEpisode))).Methods("POST")
	v1.HandleFunc("/podcast/{podcastID}/episode/{episodeID}", withAuth(makeHTTPHandleFunc(s.handleGetEpisode))).Methods("GET")
	v1.HandleFunc("/podcast/{podcastID}/episode/{episodeID}", withAuth(makeHTTPHandleFunc(s.handleUpdateEpisode))).Methods("PUT")
	v1.HandleFunc("/podcast/{podcastID}/episode/{episodeID}", withAuth(makeHTTPHandleFunc(s.handleDeleteEpisode))).Methods("DELETE")

	return http.ListenAndServe(s.listenAddr, router)
}
