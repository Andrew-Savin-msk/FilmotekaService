package filmsservice

import (
	"context"
	"errors"

	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/service"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
	"github.com/sirupsen/logrus"
)

type Films struct {
	films      store.Films
	actorFilms store.ActorFilms

	logger *logrus.Entry
	ctx    context.Context
}

func New(ctx context.Context, logger *logrus.Entry, films store.Films, actorFilms store.ActorFilms) service.Films {
	return &Films{
		films:      films,
		actorFilms: actorFilms,
		logger: logger.WithFields(logrus.Fields{
			"layer":     "service",
			"structure": "films",
		}),
		ctx: ctx,
	}
}

// TODO: Add transactions
func (f *Films) Create(film *film.Film, actorsId []int64) (int64, error) {
	id, err := f.films.Create(f.ctx, film)
	if err != nil {
		f.logger.Errorf("op: Create, unexpected error: %s", err)
		return -1, err
	}

	err = f.actorFilms.Connect(f.ctx, id, actorsId)
	if err != nil {
		f.logger.Errorf("op: Connect, unexpected error: %s", err)
		return -1, err
	}

	return id, nil
}

func (f *Films) Delete(id int64) error {
	err := f.films.Delete(f.ctx, id)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			return service.ErrNoSuchEntity
		}
		f.logger.Errorf("op: Delete, unexpected error: %s", err)
		return err
	}
	return nil
}

func (f *Films) Overwrite(film *film.Film) error {
	err := f.films.Overwrite(f.ctx, film)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			return service.ErrNoSuchEntity
		}
		f.logger.Errorf("op: Overwrite, unexpected error: %s", err)
		return err
	}
	return nil
}

func (f *Films) FindByNamePart(limit, offset int64, namePart string) ([]*film.Film, error) {
	films, err := f.films.FindMatching(f.ctx, limit, offset, namePart)
	if err != nil {
		if err != nil {
			if errors.Is(err, store.ErrRecordNotFound) {
				return nil, service.ErrNoSuchEntity
			}
			f.logger.Errorf("op: FindMatching, unexpected error: %s", err)
			return nil, err
		}
	}
	return films, nil
}

func (f *Films) GetSortedBy(limit, offset int64, sortParameter string) ([]*film.Film, error) {
	films, err := f.films.GetSorted(f.ctx, limit, offset, sortParameter)
	if err != nil {
		if errors.Is(err, store.ErrInvalidParam) {
			return nil, service.ErrInvalidParam
		}
		f.logger.Errorf("op: GetSorted, unexpected error: %s", err)
		return nil, err
	}
	return films, nil
}
