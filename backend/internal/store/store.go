package store

import (
	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
)

type Users interface {
	Create(u *user.User) error
	Find(email string) (*user.User, error)
}

type Actors interface {
	Create(actor *actor.Actor) (int64, error)
	Get(id int64) (*actor.Actor, error)
	Delete(id int64) error
	Overwrite(actor *actor.Actor) error
	GetPartitioned(limit, offset int64) ([]*actor.Actor, error)
}

type ActorFilms interface {
	GetActorsFilms(actorIds []int64) (map[int64][]*film.Film, error)
	Connect(filmId int64, actorIds []int64) error
}

type Films interface {
	Create(film *film.Film) (int64, error)
	Delete(id int64) error
	Overwrite(film *film.Film) error
	FindMatching(limit, offset int64, namePart string) ([]*film.Film, error)
	GetSorted(limit, offset int64, sortParameter string) ([]*film.Film, error)
}
