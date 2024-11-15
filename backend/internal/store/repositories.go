package store

import (
	"database/sql"

	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
)

type UserRepository interface {
	Create(u *user.User) error
	Find(id int) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}

type ActorRepository interface {
	Create(*actor.Actor) error
	Find(int) (*actor.Actor, error)
	Delete(int) (int, error)
	Overwrite(*actor.Actor) error
	GetActorsWithFilms(limit, offset int64) (map[*actor.Actor][]*film.Film, error)
}

type FilmRepository interface {
	CreateAndConnectActors(film *film.Film, actors []int) error
	Delete(id int) (int, error)
	Overwrite(film *film.Film) error
	FindByNamePart(limit, offset int64, namePart string) ([]*film.Film, error)
	FindAndSort(limit, offset int64, field string) ([]*film.Film, error)
}

type FilmActorRepository interface {
	CreateConnections(tx *sql.Tx, actors []int, filmId int) error
	GetActorsFilms(id int) ([]*film.Film, error)
}

// TODO: FilmsActors repository
