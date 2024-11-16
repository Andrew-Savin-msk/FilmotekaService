package filmsservice

import (
	"context"

	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/service"
	"github.com/sirupsen/logrus"
)

// TODO:
type Films struct {
	// filmactors repo/table connection
	// films repo/table connection

	logger *logrus.Entry
	ctx    context.Context
}

func New(ctx context.Context, logger *logrus.Entry) service.Films {
	return &Films{
		logger: logger,
		ctx:    ctx,
	}
}

func (f *Films) Create(film *film.Film, actorsId []int64) (int64, error) {
	// db Create method

	// db ConnectFilmAndActors method
	panic("unimplemented")
}

func (f *Films) Delete(id int64) error {
	// db Delete method
	panic("unimplemented")
}

func (f *Films) Overwrite(film *film.Film) error {
	// db Overwrite method
	panic("unimplemented")
}

func (f *Films) FindByNamePart(limit, offset int64, namePart string) ([]*film.Film, error) {
	// db FindByNamePart
	panic("unimplemented")
}

func (f *Films) GetSortedBy(limit, offset int64, sortParameter string) ([]*film.Film, error) {
	// db GetSortedBy method
	panic("unimplemented")
}
