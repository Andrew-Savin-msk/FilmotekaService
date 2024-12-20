package apiserver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	model "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model"
	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/repostore"
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/gorilla/mux"
)

const (
	ctxUserKey = iota
	ctxRequestKey
)

var (
	sessionName = "user-key"
)

// handleCreateUser creates a new user.
func (s *server) handleCreateUser() http.HandlerFunc {

	type request struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		u := &user.User{
			Email:   req.Email,
			Passwd:  req.Password,
			IsAdmin: false,
		}

		err = u.Validate()
		if err != nil {
			s.errorResponse(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		err = s.bc.SendEMailAddreas(u.Email)
		if err != nil {
			s.errorResponse(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		err = s.store.User().Create(u)
		if err != nil {
			s.errorResponse(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

// handleGetSession handles user login and creates a session.
func (s *server) handleGetSession() http.HandlerFunc {
	type request struct {
		Email  string `json:"email"`
		Passwd string `json:"password"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		u, err := s.store.User().FindByEmail(req.Email)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		err = u.CompareHashAndPassword(req.Passwd)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, errIncorrectEmailOrPassword)
			return
		}

		session, err := s.sessionStore.Get(r, sessionName)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		session.Values["user_id"] = u.Id

		err = session.Save(r, w)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}

}

// handleWhoamI returns information about the currently authenticated user.
func (s *server) handleWhoamI() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		u, ok := r.Context().Value(ctxUserKey).(*user.User)
		if !ok {
			s.errorResponse(w, r, http.StatusUnprocessableEntity, errNotAuthenticated)
			return
		}

		s.respond(w, r, http.StatusOK, u)
	}
}

// handleCreateActor creates a new actor.
func (s *server) handleCreateActor() http.HandlerFunc {
	type request struct {
		Name      string `json:"name"`
		Gen       string `json:"gender"`
		Birthdate string `json:"birthdate"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		birth, err := time.Parse("01-02-2006", req.Birthdate)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		act := &actor.Actor{
			Name:      req.Name,
			Gen:       req.Gen,
			Birthdate: birth,
		}

		err = act.Validate()
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		err = s.store.Actor().Create(act)
		if err != nil {
			s.errorResponse(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, act.Id)
	}
}

// handleGetActor returns information about a specific actor by ID.
func (s *server) handleGetActor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		sActorId, ok := mux.Vars(r)["actorId"]
		var actorId int64 = -1
		var err error
		if len(sActorId) != 0 || !ok {
			actorId, err = strconv.ParseInt(sActorId, 10, 32)
			if err != nil || actorId < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		act, err := s.store.Actor().Find(actorId)
		if err != nil {
			s.errorResponse(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, act)
	})
}

// handleDeleteActor deletes an actor by ID.
func (s *server) handleDeleteActor() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}
		sActorId, ok := mux.Vars(r)["actorId"]
		var actorId int64 = -1
		var err error
		if len(sActorId) != 0 || !ok {
			actorId, err = strconv.ParseInt(sActorId, 10, 32)
			if err != nil || actorId < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		id, err := s.store.Actor().Delete(actorId)
		if err != nil {
			if err == repostore.ErrRecordNotFound {
				s.respond(w, r, http.StatusOK, id)
				return
			}
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, id)
	})
}

// handleOverwriteActor updates information about an existing actor by ID.
func (s *server) handleOverwriteActor() http.Handler {
	type request struct {
		Name      string `json:"name"`
		Gen       string `json:"gender"`
		Birthdate string `json:"birthdate"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		sActorId, ok := mux.Vars(r)["actorId"]
		var actorId int64 = -1
		if len(sActorId) != 0 || !ok {
			actorId, err = strconv.ParseInt(sActorId, 10, 32)
			if err != nil || actorId < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		act := &actor.Actor{
			Id:   actorId,
			Name: req.Name,
			Gen:  req.Gen,
		}

		if req.Birthdate == "" {
			if r.Method == http.MethodPut {
				s.errorResponse(w, r, http.StatusBadRequest, fmt.Errorf("invalid date"))
				return
			}
		} else {
			birth, err := time.Parse("01-02-2006", req.Birthdate)
			if err != nil {
				s.errorResponse(w, r, http.StatusBadRequest, err)
				return
			}

			act.Birthdate = birth

			err = validation.ValidateStruct(
				act,
				validation.Field(&act.Id, validation.Required),
				validation.Field(&act.Birthdate, validation.By(model.IsDateValid())),
			)
			if err != nil {
				s.errorResponse(w, r, http.StatusBadRequest, err)
				return
			}
		}

		if r.Method == http.MethodPut {
			err = act.Validate()
			if err != nil {
				s.errorResponse(w, r, http.StatusBadRequest, err)
				return
			}
		}

		err = s.store.Actor().Overwrite(act)
		if err != nil {
			if err == repostore.ErrRecordNotFound {
				s.respond(w, r, http.StatusOK, "")
				return
			}
			s.errorResponse(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, "")
	})
}

// handleGetActors returns a list of all actors with their films.
func (s *server) handleGetActors() http.Handler {
	type respond struct {
		Act   *actor.Actor `json:"actor"`
		Films []*film.Film `json:"films"`
	}

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		var limit int64 = 5
		var err error
		if len(limitStr) != 0 {
			limit, err = strconv.ParseInt(limitStr, 10, 32)
			if err != nil || limit < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		offsetStr := r.URL.Query().Get("offset")
		var offset int64 = 5
		if len(offsetStr) != 0 {
			offset, err = strconv.ParseInt(offsetStr, 10, 32)
			if err != nil || offset < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		actors, err := s.store.Actor().GetActorsWithFilms(limit, offset)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		res := []respond{}
		for actor, films := range actors {
			res = append(res, respond{Act: actor, Films: films})
		}

		s.respond(w, r, http.StatusOK, res)
	})
}

// handleCreateFilm creates a new film.
func (s *server) handleCreateFilm() http.Handler {
	type request struct {
		Id        int     `json:"id"`
		Name      string  `json:"name"`
		Desc      string  `json:"description,omitempty"`
		Date      string  `json:"release_date"`
		Assesment float32 `json:"assesment"`
		Actors    []int   `json:"actors"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		req := &request{}

		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		date, err := time.Parse("01-02-2006", req.Date)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		film := &film.Film{
			Name:      req.Name,
			Desc:      req.Desc,
			Date:      date,
			Assesment: req.Assesment,
		}

		fmt.Println(req, film)

		err = s.store.Film().CreateAndConnectActors(film, req.Actors)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, film.Id)
	})
}

// handleDeleteFilm deletes a film by ID.
func (s *server) handleDeleteFilm() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		sFilmId, ok := mux.Vars(r)["filmId"]
		var filmId int64 = -1
		var err error
		if len(sFilmId) != 0 || !ok {
			filmId, err = strconv.ParseInt(sFilmId, 10, 32)
			if err != nil || filmId < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		id, err := s.store.Film().Delete(filmId)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, id)
	})
}

// handleOverwriteFilm updates information about an existing film by ID.
func (s *server) handleOverwriteFilm() http.Handler {
	type request struct {
		Name      string  `json:"name"`
		Desc      string  `json:"description"`
		Date      string  `json:"release_date"`
		Assesment float32 `json:"assesment"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		sFilmId, ok := mux.Vars(r)["filmId"]
		var filmId int64 = -1
		if len(sFilmId) != 0 || !ok {
			filmId, err = strconv.ParseInt(sFilmId, 10, 32)
			if err != nil || filmId < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		film := &film.Film{
			Id:        filmId,
			Name:      req.Name,
			Desc:      req.Desc,
			Assesment: req.Assesment,
		}

		if req.Date == "" {
			if r.Method == http.MethodPut {
				s.errorResponse(w, r, http.StatusBadRequest, fmt.Errorf("invalid date"))
				return
			}
		} else {
			birth, err := time.Parse("01-02-2006", req.Date)
			if err != nil {
				s.errorResponse(w, r, http.StatusBadRequest, err)
				return
			}

			film.Date = birth

			err = validation.ValidateStruct(
				film,
				validation.Field(&film.Id, validation.Required),
				validation.Field(&film.Date, validation.By(model.IsDateValid())),
			)
			if err != nil {
				s.errorResponse(w, r, http.StatusBadRequest, err)
				return
			}
		}

		if r.Method == http.MethodPut {
			err = film.Validate()
			if err != nil {
				s.errorResponse(w, r, http.StatusBadRequest, err)
				return
			}
		}

		err = s.store.Film().Overwrite(film)
		if err != nil {
			if err == repostore.ErrRecordNotFound {
				s.respond(w, r, http.StatusOK, "")
				return
			}
			s.errorResponse(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, "")
	})
}

// handleFindFilmByNamePart finds a film by part of its name.
func (s *server) handleFindFilmByNamePart() http.Handler {
	type request struct {
		NamePart string `json:"name_part"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		var limit int64 = 5
		if len(limitStr) != 0 {
			limit, err = strconv.ParseInt(limitStr, 10, 32)
			if err != nil || limit < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		offsetStr := r.URL.Query().Get("offset")
		var offset int64 = 5
		if len(offsetStr) != 0 {
			offset, err = strconv.ParseInt(offsetStr, 10, 32)
			if err != nil || offset < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		films, err := s.store.Film().FindByNamePart(limit, offset, req.NamePart)
		if err != nil {
			if err == repostore.ErrRecordNotFound {
				s.respond(w, r, http.StatusOK, struct{}{})
				return
			}
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, films)
	})
}

// handleGetSortedFilms returns a list of films sorted by a specific criterion.
func (s *server) handleGetSortedFilms() http.Handler {
	type request struct {
		SortParam string `json:"sorting_parameter"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		limitStr := r.URL.Query().Get("limit")
		var limit int64 = 5
		if len(limitStr) != 0 {
			limit, err = strconv.ParseInt(limitStr, 10, 32)
			if err != nil || limit < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		offsetStr := r.URL.Query().Get("offset")
		var offset int64 = 5
		if len(offsetStr) != 0 {
			offset, err = strconv.ParseInt(offsetStr, 10, 32)
			if err != nil || offset < 0 {
				s.errorResponse(w, r, http.StatusBadRequest, ErrInvalidQuerryParams)
				return
			}
		}

		films, err := s.store.Film().FindAndSort(limit, offset, req.SortParam)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, films)
	})

}

// Interface methods

func (s *server) ServerHTTP(w http.ResponseWriter, r *http.Request) {
	s.mux.ServeHTTP(w, r)
}
