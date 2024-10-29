package apiserver

import (
	"context"
	"fmt"
	"net/http"
	"time"

	model "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

// Path middlewares

func (s *server) basePaths(next http.Handler) http.Handler {
	return s.wrapSetRequestId(s.wrapLogRequest(next))
}

func (s *server) protectedPaths(next http.Handler) http.Handler {
	return s.basePaths(s.wrapAuthorise(next))
}

func (s *server) adminPaths(next http.Handler) http.Handler {
	return s.protectedPaths(s.wrapAdminCheck(next))
}

// Middleware wrappers

func (s *server) wrapSetRequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-Id", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxRequestKey, id)))
	})
}

func (s *server) wrapLogRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxRequestKey),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)
		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}

		next.ServeHTTP(rw, r)

		logger.Infof(
			"completed with %v %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Since(start),
		)
	})
}

func (s *server) wrapAuthorise(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Getting cookies
		cookie, err := r.Cookie(sessionName)
		fmt.Println(cookie, cookie.Name)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		// Geting session data
		session, err := s.sessionStore.Get(r, cookie.Name)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		// Geting session token
		id, ok := session.Values["user_id"].(int)
		if !ok {
			s.errorResponse(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		u, err := s.store.User().Find(id)
		if err != nil {
			s.errorResponse(w, r, http.StatusUnauthorized, errNotAuthenticated)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxUserKey, u)))
	})
}

func (s *server) wrapAdminCheck(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := r.Context().Value(ctxUserKey).(*model.User)
		s.logger.Debugf("Authorised user: %v", user)
		if !ok {
			s.errorResponse(w, r, http.StatusUnauthorized, errResourceForbiden)
			return
		}

		if !user.IsAdmin {
			s.errorResponse(w, r, http.StatusForbidden, errResourceForbiden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
