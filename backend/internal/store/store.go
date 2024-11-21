package store

import (
	"context"

	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
)

type Users interface {
	Create(ctx context.Context, u *user.User) error
	Find(ctx context.Context, email string) (*user.User, error)
}

type Actors interface {
	Create(ctx context.Context, actor *actor.Actor) (int64, error)
	Get(ctx context.Context, id int64) (*actor.Actor, error)
	Delete(ctx context.Context, id int64) error
	Overwrite(ctx context.Context, actor *actor.Actor) error
	GetPartitioned(ctx context.Context, limit, offset int64) ([]*actor.Actor, error)
}

type ActorFilms interface {
	GetActorsFilms(ctx context.Context, actorIds []int64) (map[int64][]*film.Film, error)
	Connect(ctx context.Context, filmId int64, actorIds []int64) error
}

type Films interface {
	Create(ctx context.Context, film *film.Film) (int64, error)
	Delete(ctx context.Context, id int64) error
	Overwrite(ctx context.Context, film *film.Film) error
	FindMatching(ctx context.Context, limit, offset int64, namePart string) ([]*film.Film, error)
	GetSorted(ctx context.Context, limit, offset int64, sortParameter string) ([]*film.Film, error)
}
