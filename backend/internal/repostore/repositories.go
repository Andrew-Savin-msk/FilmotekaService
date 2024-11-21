package repostore

import (
	"database/sql"

	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
)

type UserRepository interface {
	Create(u *user.User) error
	Find(id int64) (*user.User, error)
	FindByEmail(email string) (*user.User, error)
}

type ActorRepository interface {
	Create(*actor.Actor) error
	Find(int64) (*actor.Actor, error)
	Delete(int64) (int64, error)
	Overwrite(*actor.Actor) error
	GetActorsWithFilms(limit, offset int64) (map[*actor.Actor][]*film.Film, error)
}

type FilmRepository interface {
	CreateAndConnectActors(film *film.Film, actors []int) error
	Delete(id int64) (int64, error)
	Overwrite(film *film.Film) error
	FindByNamePart(limit, offset int64, namePart string) ([]*film.Film, error)
	FindAndSort(limit, offset int64, field string) ([]*film.Film, error)
}

type FilmActorRepository interface {
	CreateConnections(tx *sql.Tx, actors []int, filmId int64) error
	GetActorsFilms(id int64) ([]*film.Film, error)
}

// TODO: FilmsActors repository
