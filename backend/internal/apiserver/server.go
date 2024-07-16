package apiserver

import (
	"net/http"

	brockerclient "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/broker_client"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/config"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type server struct {
	mux          *http.ServeMux
	sessionStore sessions.Store
	store        store.Store
	bc           brockerclient.Client
	logger       *logrus.Logger
}

func newServer(st store.Store, bc brockerclient.Client, logger *logrus.Logger, cfg *config.Config) *server {
	srv := &server{
		mux:          http.NewServeMux(),
		logger:       logger,
		store:        st,
		sessionStore: sessions.NewCookieStore([]byte(cfg.Srv.SessionKey)),
		bc:           bc,
	}

	srv.setMuxer()

	return srv
}

func (s *server) setMuxer() {
	// Public endpoints
	s.mux.Handle("/register", s.basePaths(s.handleCreateUser()))
	s.mux.Handle("/authorize", s.basePaths(s.handleGetSession()))
	s.mux.Handle("/get-actor", s.basePaths(s.handleGetActor()))
	s.mux.Handle("/get-actors", s.basePaths(s.handleGetActors()))
	s.mux.Handle("/films", s.basePaths(s.handleFindFilmByNamePart()))
	s.mux.Handle("/select-films", s.basePaths(s.handleGetSortedFilms()))

	// Authorisation required endpoints
	s.mux.Handle("/private/who-am-i", s.protectedPaths(s.handleWhoamI()))

	// Admin rights required endpoints
	s.mux.Handle("/private/create-actor", s.adminPaths(s.handleCreateActor()))
	s.mux.Handle("/private/delete-actor", s.adminPaths(s.handleDeleteActor()))
	s.mux.Handle("/private/update-actor", s.adminPaths(s.handleOverwrightActor()))
	s.mux.Handle("/private/post-film", s.adminPaths(s.handleCreateFilm()))
	s.mux.Handle("/private/delete-film", s.adminPaths(s.handleDeleteFilm()))
	s.mux.Handle("/private/update-film", s.adminPaths(s.handleOverwrightFilm()))
}
