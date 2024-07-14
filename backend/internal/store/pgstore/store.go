package pgstore

import (
	"database/sql"

	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
)

type Store struct {
	db  *sql.DB
	ur  *UserRepository
	ar  *ActorRepository
	fr  *FilmRepository
	far *FilmActorRepository
}

func New(db *sql.DB) *Store {
	return &Store{
		db: db,
	}
}

func (s Store) Close() {
	s.db.Close()
}

func (s *Store) User() store.UserRepository {
	if s.ur == nil {
		s.ur = &UserRepository{
			st: s,
		}
	}
	return s.ur
}

func (s *Store) Actor() store.ActorRepository {
	if s.ar == nil {
		s.ar = &ActorRepository{
			st: s,
		}
	}
	return s.ar
}

func (s *Store) Film() store.FilmRepository {
	if s.fr == nil {
		s.fr = &FilmRepository{
			st: s,
		}
	}
	return s.fr
}

func (s *Store) FilmActor() store.FilmActorRepository {
	if s.far == nil {
		s.far = &FilmActorRepository{
			st: s,
		}
	}
	return s.far
}
