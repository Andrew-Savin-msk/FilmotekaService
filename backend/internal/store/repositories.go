package store

import (
	"database/sql"

	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
)

type UserRepository interface {
	Create(*user.User) error
	Find(int) (*user.User, error)
	FindByEmail(string) (*user.User, error)
}

type ActorRepository interface {
	Create(*actor.Actor) error
	Find(int) (*actor.Actor, error)
	Delete(int) (int, error)
	Overwright(*actor.Actor) error
	GetActorsWithFilms() (map[*actor.Actor][]*film.Film, error)
}

type FilmRepository interface {
	CreateAndConnectActors(*film.Film, []int) error
	Delete(int) (int, error)
	Overwright(*film.Film) error
	FindByNamePart(string) (*film.Film, error)
	FindAndSort(string) ([]*film.Film, error)
}

type FilmActorRepository interface {
	CreateConnections(*sql.Tx, []int, int) error
	GetActorsFilms(int) ([]*film.Film, error)
}

// TODO: FilmsActors repository
