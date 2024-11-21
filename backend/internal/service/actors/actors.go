package actorsservice

import (
	"context"
	"errors"

	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/service"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
	"github.com/sirupsen/logrus"
)

type Actors struct {
	actors     store.Actors
	actorFilms store.ActorFilms

	logger *logrus.Entry
	ctx    context.Context
}

func New(ctx context.Context, logger *logrus.Entry, actors store.Actors, actorFilms store.ActorFilms) service.Actors {
	return &Actors{
		actors:     actors,
		actorFilms: actorFilms,
		logger: logger.WithFields(logrus.Fields{
			"layer":     "service",
			"structure": "actors",
		}),
		ctx: ctx,
	}
}

func (ac *Actors) Create(actor *actor.Actor) (int64, error) {
	id, err := ac.actors.Create(ac.ctx, actor)
	if err != nil {
		ac.logger.Errorf("op: Create, unexpected error: %s", err)
		return -1, err
	}

	return id, nil
}

func (ac *Actors) GetById(id int64) (*actor.Actor, error) {
	a, err := ac.actors.Get(ac.ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, service.ErrNoSuchEntity
		}
		ac.logger.Errorf("op: GetById, unexpected error: %s", err)
		return nil, err
	}
	return a, nil
}

func (ac *Actors) Delete(id int64) error {
	err := ac.actors.Delete(ac.ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			return service.ErrNoSuchEntity
		}
		ac.logger.Errorf("op: Delete, unexpected error: %s", err)
		return err
	}
	return nil
}

func (ac *Actors) Overwrite(actor *actor.Actor) error {
	err := ac.actors.Overwrite(ac.ctx, actor)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			return service.ErrNoSuchEntity
		}
		ac.logger.Errorf("op: Overwrite, unexpected error: %s", err)
		return err
	}
	return nil
}

func (ac *Actors) GetInfo(limit, offset int64) (map[*actor.Actor][]*film.Film, error) {
	actors, err := ac.actors.GetPartitioned(ac.ctx, limit, offset)
	if err != nil {
		ac.logger.Errorf("op: GetPartitioned, unexpected error: %s", err)
		return nil, err
	}

	res := map[*actor.Actor][]*film.Film{}
	actorIds := make([]int64, len(actors))
	for i, actor := range actors {
		actorIds[i] = actor.Id
		res[actor] = []*film.Film{}
	}

	films, err := ac.actorFilms.GetActorsFilms(ac.ctx, actorIds)
	if err != nil {
		ac.logger.Errorf("op: GetActorsFilms, unexpected error: %s", err)
		return nil, err
	}

	for actor, _ := range res {
		res[actor] = films[actor.Id]
	}

	return res, nil
}
