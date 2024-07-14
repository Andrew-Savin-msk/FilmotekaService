package apiserver

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	model "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model"
	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
	validation "github.com/go-ozzo/ozzo-validation"
)

const (
	ctxUserKey = iota
	ctxRequestKey
)

var (
	sessionName = "user-key"
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not auntificated")
	errResourceForbiden         = errors.New("you dont have permossions to get this resource")
	errIncorrectId              = errors.New("presented incorrect id type")
	errMethodNotAllowed         = errors.New("unsuportable method type")
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
	type request struct {
		Id int `json:"id"`
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

		act, err := s.store.Actor().Find(req.Id)
		if err != nil {
			s.errorResponse(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusOK, act)
	})
}

// handleDeleteActor deletes an actor by ID.
func (s *server) handleDeleteActor() http.Handler {
	type request struct {
		Id int `json:"id"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}
		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		id, err := s.store.Actor().Delete(req.Id)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.respond(w, r, http.StatusOK, id)
			}
			s.errorResponse(w, r, http.StatusBadRequest, err)
		}

		s.respond(w, r, http.StatusOK, id)
	})
}

// handleOverwrightActor updates information about an existing actor by ID.
func (s *server) handleOverwrightActor() http.Handler {
	type request struct {
		Id        int    `json:"id"`
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

		act := &actor.Actor{
			Id:   req.Id,
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

		err = s.store.Actor().Overwright(act)
		if err != nil {
			if err == store.ErrRecordNotFound {
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

		actors, err := s.store.Actor().GetActorsWithFilms()
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		res := []respond{}
		for actor, films := range actors {
			res = append(res, respond{Act: actor, Films: films})
			fmt.Println(actor, films)
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
	type request struct {
		Id int `json:"id"`
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			s.errorResponse(w, r, http.StatusMethodNotAllowed, errMethodNotAllowed)
			return
		}

		req := &request{}
		err := json.NewDecoder(r.Body).Decode(req)
		if err != nil {
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		id, err := s.store.Film().Delete(req.Id)
		if err != nil {
			s.errorResponse(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, id)
	})
}

// handleOverwrightFilm updates information about an existing film by ID.
func (s *server) handleOverwrightFilm() http.Handler {
	type request struct {
		Id        int     `json:"id"`
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

		film := &film.Film{
			Id:        req.Id,
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

		err = s.store.Film().Overwright(film)
		if err != nil {
			if err == store.ErrRecordNotFound {
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

		film, err := s.store.Film().FindByNamePart(req.NamePart)
		if err != nil {
			if err == store.ErrRecordNotFound {
				s.respond(w, r, http.StatusOK, struct{}{})
				return
			}
			s.errorResponse(w, r, http.StatusBadRequest, err)
			return
		}

		s.respond(w, r, http.StatusOK, film)
	})
}

// handleGetSortedFilms returns a list of films sorted by a specific criterion.
func (s *server) handleGetSortedFilms() http.Handler {
	type request struct {
		SortParam string `json:"sorting_parameter"`
		Amount    int    `json:"amount"`
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

		films, err := s.store.Film().FindAndSort(req.SortParam, req.Amount)
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
