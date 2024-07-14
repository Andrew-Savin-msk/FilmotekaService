package pgstore

import (
	"database/sql"
	"fmt"

	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
)

type FilmRepository struct {
	st *Store
}

func (f *FilmRepository) CreateAndConnectActors(film *film.Film, actors []int) error {

	tx, err := f.st.db.Begin()
	if err != nil {
		// TODO: Think about special Err statment to detect internal issues during server running
		return err
	}

	err = tx.QueryRow(
		"INSERT INTO films (name, description, release_date, assesment) VALUES ($1, $2, $3, $4) RETURNING id",
		film.Name,
		film.Desc,
		film.Date,
		film.Assesment,
	).Scan(&film.Id)
	if err != nil {
		tErr := tx.Rollback()
		if tErr != nil {
			return fmt.Errorf("rollback error: %e, triggered by error: %e", tErr, err)
		}
		return err
	}

	err = f.st.FilmActor().CreateConnections(tx, actors, film.Id)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (f *FilmRepository) Delete(id int) (int, error) {
	res, err := f.st.db.Exec(
		"DELETE FROM films WHERE id = $1",
		id,
	)
	if err != nil {
		return -1, err
	}

	am, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	if am == 0 {
		return -1, store.ErrRecordNotFound
	}

	return int(am), nil
}

func (f *FilmRepository) Overwright(film *film.Film) error {
	res, err := f.st.db.Exec(
		"UPDATE films SET "+
			"name = CASE WHEN name <> $1 AND $1 <> '' THEN $1 ELSE name END, "+
			"description = CASE WHEN description <> $2 AND $2 <> '' THEN $2 ELSE description END, "+
			"release_date = CASE WHEN release_date <> $3 AND $4 IS NOT FALSE THEN $3 ELSE release_date END, "+
			"assesment = CASE WHEN assesment <> $5 AND assesment BETWEEN 0 AND 10 THEN $5 ELSE assesment END "+
			"WHERE id = $6",
		film.Name,
		film.Desc,
		film.Date,
		film.Date.IsZero(),
		film.Assesment,
		film.Id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return store.ErrRecordNotFound
	}

	return nil
}

func (f *FilmRepository) FindByNamePart(namePart string) (*film.Film, error) {
	film := &film.Film{}
	err := f.st.db.QueryRow(
		"SELECT id, name, description, release_date, assesment FROM films "+
			"WHERE name LIKE '%' || $1 || '%'",
		namePart,
	).Scan(
		&film.Id,
		&film.Name,
		&film.Desc,
		&film.Date,
		&film.Assesment,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return film, store.ErrRecordNotFound
		}
		return nil, err
	}
	return film, nil
}

func (f *FilmRepository) FindAndSort(field string, amount int) ([]*film.Film, error) {
	films := []*film.Film{}
	if field != "name" && field != "release_date" && field != "assesment" {
		return films, store.ErrForbiddenParameters
	}

	rows, err := f.st.db.Query(
		"SELECT id, name, description, release_date, assesment FROM films "+
			"ORDER BY $1 DESC "+
			"LIMIT $2",
		field,
		amount,
	)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var film film.Film
		err = rows.Scan(&film.Id, &film.Name, &film.Desc, &film.Date, &film.Assesment)
		if err != nil {
			return nil, err
		}
		films = append(films, &film)
	}

	return films, nil
}
