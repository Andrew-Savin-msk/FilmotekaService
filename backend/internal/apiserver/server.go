package apiserver

import (
	brockerclient "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/broker_client"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/config"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/sirupsen/logrus"
)

type server struct {
	mux          *mux.Router
	sessionStore sessions.Store
	store        store.Store
	bc           brockerclient.Client
	logger       *logrus.Logger
}

func newServer(st store.Store, bc brockerclient.Client, logger *logrus.Logger, cfg *config.Config) *server {
	srv := &server{
		mux:          mux.NewRouter(),
		logger:       logger,
		store:        st,
		sessionStore: sessions.NewCookieStore([]byte(cfg.Srv.SessionKey)),
		bc:           bc,
	}

	srv.setMuxer()

	return srv
}

func (s *server) setMuxer() {

	s.mux.Use(s.wrapSetRequestId)
	s.mux.Use(s.wrapLogRequest)

	s.mux.Handle("/register", s.handleCreateUser())
	s.mux.Handle("/authorize", s.handleGetSession())
	s.mux.Handle("/get-actor/{actorId}", s.handleGetActor())
	s.mux.Handle("/get-actors", s.handleGetActors())        // TODO: Pagination (offset, limit)
	s.mux.Handle("/films", s.handleFindFilmByNamePart())    // TODO: Pagination (offset, limit)
	s.mux.Handle("/select-films", s.handleGetSortedFilms()) // TODO: Pagination (offset, limit)

	// Authorisation required endpoints
	private := s.mux.NewRoute().Subrouter()
	private.Use(s.wrapAuthorise)
	private.Handle("/private/who-am-i", s.handleWhoamI())

	admin := private.NewRoute().Subrouter()
	admin.Use(s.wrapAdminCheck)
	// Admin rights required endpoints
	admin.Handle("/private/create-actor", s.handleCreateActor())
	admin.Handle("/private/delete-actor/{actorId}", s.handleDeleteActor())
	admin.Handle("/private/update-actor{actorId}", s.handleOverwriteActor())
	admin.Handle("/private/post-film", s.handleCreateFilm())
	admin.Handle("/private/delete-film/{filmId}", s.handleDeleteFilm())
	admin.Handle("/private/update-film/{filmId}", s.handleOverwriteFilm())
}
