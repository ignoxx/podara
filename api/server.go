package api

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/ignoxx/podara/poc3/storage"
	"github.com/rs/cors"
)

type Server struct {
	listenAddr string
	store      storage.Storage
	imageDir   string
	audioDir   string
}

func NewServer(listenAddr string, store storage.Storage, imageDir, audioDir string) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
		imageDir:   imageDir,
		audioDir:   audioDir,
	}
}

func (s *Server) Start() error {
	router := mux.NewRouter()

	router.HandleFunc("/", s.handleLandingPage)
	router.HandleFunc("/login", s.handleLoginPage)
	router.HandleFunc("/register", s.handleRegisterPage)
    router.HandleFunc("/profile", s.handleProfilePage)

	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/register", makeHTTPHandleFunc(s.handleCreateUser)).Methods("POST")
	v1.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods("POST")

	v1.HandleFunc("/podcast", withAuth(makeHTTPHandleFunc(s.handleCreatePodcast))).Methods("POST")
	v1.HandleFunc("/podcast", withAuth(makeHTTPHandleFunc(s.handleGetAllPodcasts))).Methods("GET")
	v1.HandleFunc("/podcast/{podcast_id}", withAuth(makeHTTPHandleFunc(s.handleGetPodcast))).Methods("GET")
	v1.HandleFunc("/podcast/{podcast_id}", withAuth(makeHTTPHandleFunc(s.handleUpdatePodcast))).Methods("PUT")
	v1.HandleFunc("/podcast/{podcast_id}", withAuth(makeHTTPHandleFunc(s.handleDeletePodcast))).Methods("DELETE")

	v1.HandleFunc("/podcast/{podcast_id}/rss.xml", makeHTTPHandleFunc(s.handleGetPodcastRss)).Methods("GET", "HEAD")

	v1.HandleFunc("/podcast/{podcast_id}/episode", withAuth(makeHTTPHandleFunc(s.handleCreateEpisode))).Methods("POST")
	v1.HandleFunc("/podcast/{podcast_id}/episodes", withAuth(makeHTTPHandleFunc(s.handleGetEpisodes))).Methods("GET")
	v1.HandleFunc("/podcast/{podcast_id}/episode/{episode_id}", withAuth(makeHTTPHandleFunc(s.handleGetEpisode))).Methods("GET")
	v1.HandleFunc("/podcast/{podcast_id}/episode/{episode_id}", withAuth(makeHTTPHandleFunc(s.handleUpdateEpisode))).Methods("PUT")
	v1.HandleFunc("/podcast/{podcast_id}/episode/{episode_id}", withAuth(makeHTTPHandleFunc(s.handleDeleteEpisode))).Methods("DELETE")

	handler := cors.AllowAll().Handler(router)

	return http.ListenAndServe(s.listenAddr, handler)
}

func (s *Server) handleLandingPage(w http.ResponseWriter, r *http.Request) {
	files := []string{"templates/base.tmpl", "templates/index.tmpl"}
	templates := template.Must(template.ParseFiles(files...))

	userClaims, _ := getJwtClaims(r)

	data := struct {
		User *UserClaims
	}{User: userClaims}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleLoginPage(w http.ResponseWriter, r *http.Request) {
	files := []string{"templates/base.tmpl", "templates/login.tmpl"}
	templates := template.Must(template.ParseFiles(files...))

	userClaims, _ := getJwtClaims(r)

	data := struct {
		User *UserClaims
	}{User: userClaims}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleRegisterPage(w http.ResponseWriter, r *http.Request) {
	files := []string{"templates/base.tmpl", "templates/register.tmpl"}
	templates := template.Must(template.ParseFiles(files...))

	userClaims, _ := getJwtClaims(r)

	data := struct {
		User *UserClaims
	}{User: userClaims}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handleProfilePage(w http.ResponseWriter, r *http.Request) {
	files := []string{"templates/base.tmpl", "templates/profile.tmpl"}
	templates := template.Must(template.ParseFiles(files...))

	userClaims, _ := getJwtClaims(r)

	data := struct {
		User *UserClaims
	}{User: userClaims}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
