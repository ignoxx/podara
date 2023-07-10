package api

import (
	"net/http"
	"text/template"

	"github.com/gorilla/mux"
	"github.com/ignoxx/podara/poc3/storage"
	"github.com/ignoxx/podara/poc3/types"
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

	// serve static files
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir(s.imageDir))))
	router.PathPrefix("/audio/").Handler(http.StripPrefix("/audio/", http.FileServer(http.Dir(s.audioDir))))

	router.HandleFunc("/", s.handleLandingPage)
	router.HandleFunc("/login", s.handleLoginPage)
	router.HandleFunc("/register", s.handleRegisterPage)
	router.HandleFunc("/profile", s.handleProfilePage)
	router.HandleFunc("/podcasts", s.handlePodcastsPage)
	router.HandleFunc("/podcast/{podcast_id}", s.handlePodcastPage)

	v1 := router.PathPrefix("/api/v1").Subrouter()

	v1.HandleFunc("/register", makeHTTPHandleFunc(s.handleCreateUser)).Methods("POST")
	v1.HandleFunc("/login", makeHTTPHandleFunc(s.handleLogin)).Methods("POST")

	podcastRouter := v1.PathPrefix("/podcast").Subrouter()

	podcastRouter.HandleFunc("/", withAuth(makeHTTPHandleFunc(s.handleCreatePodcast))).Methods("POST")
	podcastRouter.HandleFunc("/", withAuth(makeHTTPHandleFunc(s.handleGetAllPodcasts))).Methods("GET")
	podcastRouter.HandleFunc("/{podcast_id}", withAuth(makeHTTPHandleFunc(s.handleGetPodcast))).Methods("GET")
	podcastRouter.HandleFunc("/{podcast_id}", withAuth(makeHTTPHandleFunc(s.handleUpdatePodcast))).Methods("PUT")
	podcastRouter.HandleFunc("/{podcast_id}", withAuth(makeHTTPHandleFunc(s.handleDeletePodcast))).Methods("DELETE")

	podcastRouter.HandleFunc("/{podcast_id}/rss.xml", makeHTTPHandleFunc(s.handleGetPodcastRss)).Methods("GET", "HEAD")

	podcastRouter.HandleFunc("/{podcast_id}/episode", withAuth(makeHTTPHandleFunc(s.handleCreateEpisode))).Methods("POST")
	podcastRouter.HandleFunc("/{podcast_id}/episodes", withAuth(makeHTTPHandleFunc(s.handleGetEpisodes))).Methods("GET")
	podcastRouter.HandleFunc("/{podcast_id}/episode/{episode_id}", withAuth(makeHTTPHandleFunc(s.handleGetEpisode))).Methods("GET")
	podcastRouter.HandleFunc("/{podcast_id}/episode/{episode_id}", withAuth(makeHTTPHandleFunc(s.handleUpdateEpisode))).Methods("PUT")
	podcastRouter.HandleFunc("/{podcast_id}/episode/{episode_id}", withAuth(makeHTTPHandleFunc(s.handleDeleteEpisode))).Methods("DELETE")

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

func (s *Server) handlePodcastsPage(w http.ResponseWriter, r *http.Request) {
	files := []string{"templates/base.tmpl", "templates/podcasts.tmpl"}
	templates := template.Must(template.ParseFiles(files...))

	userClaims, _ := getJwtClaims(r)

	podcasts, err := s.store.GetPodcastByUserID(userClaims.Id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		User     *UserClaims
		Podcasts []*types.Podcast
	}{User: userClaims, Podcasts: podcasts}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (s *Server) handlePodcastPage(w http.ResponseWriter, r *http.Request) {
	podcast_id := mux.Vars(r)["podcast_id"]

	files := []string{"templates/base.tmpl", "templates/podcast.tmpl"}
	templates := template.Must(template.ParseFiles(files...))

	userClaims, _ := getJwtClaims(r)

	episodes, err := s.store.GetAllEpisodes(podcast_id)
	podcast, err := s.store.GetPodcastByID(podcast_id)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := struct {
		User     *UserClaims
		Episodes []*types.Episode
		Podcast  *types.Podcast
	}{User: userClaims, Episodes: episodes, Podcast: podcast}

	if err := templates.ExecuteTemplate(w, "base", data); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
