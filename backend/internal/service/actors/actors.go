package actorsservice

import (
	"context"

	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/service"
	"github.com/sirupsen/logrus"
)

// TODO:
type Actors struct {
	// actors repo/table connection
	// films repo/table connection

	logger *logrus.Entry
	ctx    context.Context
}

func New(ctx context.Context, logger *logrus.Entry) service.Actors {
	return &Actors{
		logger: logger,
		ctx:    ctx,
	}
}

func (ac *Actors) Create(actor *actor.Actor) (int64, error) {
	// db Create method
	panic("unimplemented")
}

func (ac *Actors) GetById(id int64) (*actor.Actor, error) {
	// db Find method
	panic("unimplemented")
}

func (ac *Actors) Delete(id int64) error {
	// db Delete method
	panic("unimplemented")
}

func (ac *Actors) Overwrite(actor *actor.Actor) error {
	// db Overwrite method
	panic("unimplemented")
}

func (ac *Actors) GetInfo(limit, offset int64) (map[*actor.Actor][]*film.Film, error) {
	// db GetActors method

	// db GetActorsFilms method
	panic("unimplemented")
}
