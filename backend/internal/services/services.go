package services

import (
	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
)

type Users interface {
	Create(u *user.User) error
	GetByEmail(email string) (*user.User, error)
}

type Actors interface {
	Create(actor *actor.Actor) (int64, error)
	GetById(id int64) (*actor.Actor, error)
	Delete(id int64) error
	Overwrite(actor *actor.Actor) error
	GetInfo(limit, offset int64)
}

type Films interface {
	Create(film *film.Film, actorsId []int64) (int64, error)
	Delete(id int64) error
	Overwrite(film *film.Film) error
	FindByNamePart(limit, offset int64, namePart string) ([]*film.Film, error)
	FindAndSort(limit, offset int64, sortParameter string)
}
