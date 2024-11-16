package pgstore

import (
	"database/sql"
	"fmt"

	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
)

type FilmActorRepository struct {
	st *Store
}

// Doesn't forget rollback in any error case
func (far *FilmActorRepository) CreateConnections(tx *sql.Tx, actors []int, filmId int) error {
	stmt, err := tx.Prepare("INSERT INTO films_actors (film_id, actor_id) VALUES ($1, $2)")
	// No rollback due to no changes
	if err != nil {
		tErr := tx.Rollback()
		if tErr != nil {
			return fmt.Errorf("rollback error: %e, triggered by error: %e", tErr, err)
		}
		return err
	}

	for _, actorId := range actors {
		_, err := stmt.Exec(filmId, actorId)
		if err != nil {
			tErr := tx.Rollback()
			if tErr != nil {
				return fmt.Errorf("rollback error: %e, triggered by error: %e", tErr, err)
			}
			return err
		}
	}

	return nil
}

// TODO:
func (far *FilmActorRepository) GetActorsFilms(id int) ([]*film.Film, error) {
	films := []*film.Film{}
	rows, err := far.st.db.Query(
		"SELECT f.id, f.name, f.description, f.release_date, f.assesment FROM films_actors AS fa "+
			"JOIN films AS f ON fa.film_id = f.id "+
			"WHERE fa.actor_id = $1 LIMIT 5",
		id,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var film film.Film
		err := rows.Scan(&film.Id, &film.Name, &film.Desc, &film.Date, &film.Assesment)
		if err != nil {
			return nil, err
		}
		films = append(films, &film)
	}

	return films, nil
}
