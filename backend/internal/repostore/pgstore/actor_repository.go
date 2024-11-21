package pgstore

import (
	"database/sql"

	actor "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/actor"
	film "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/film"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/repostore"
)

type ActorRepository struct {
	st *Store
}

func (a *ActorRepository) Create(act *actor.Actor) error {
	err := act.Validate()
	if err != nil {
		return err
	}

	// TODO: Change to Exec statment
	return a.st.db.QueryRow(
		"INSERT INTO actors (name, gender, birthdate) VALUES ($1, $2, $3) RETURNING id",
		act.Name,
		act.Gen,
		act.Birthdate,
	).Scan(
		&act.Id,
	)
}

func (a *ActorRepository) Find(id int64) (*actor.Actor, error) {
	act := &actor.Actor{
		Id: id,
	}

	err := a.st.db.QueryRow(
		"SELECT gender, birthdate, name FROM actors WHERE id = $1",
		id,
	).Scan(
		&act.Gen,
		&act.Birthdate,
		&act.Name,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, repostore.ErrRecordNotFound
		}
		return nil, err
	}
	return act, nil
}

func (a *ActorRepository) Delete(id int64) (int64, error) {
	res, err := a.st.db.Exec(
		"DELETE FROM actors WHERE id = $1",
		id,
	)
	if err != nil {
		return -1, err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return -1, err
	}
	if rows == 0 {
		return -1, repostore.ErrRecordNotFound
	}
	return rows, nil
}

func (a *ActorRepository) Overwrite(act *actor.Actor) error {
	res, err := a.st.db.Exec(
		"UPDATE actors SET "+
			"gender = CASE WHEN gender <> $1 AND $1 <> '' THEN $1 ELSE gender END, "+
			"birthdate = CASE WHEN birthdate <> $2 AND $3 IS NOT FALSE THEN $2 ELSE birthdate END, "+
			"name = CASE WHEN name <> $4 AND $4 <> '' THEN $4 ELSE name END "+
			"WHERE id = $5",
		act.Gen,
		act.Birthdate,
		act.Birthdate.IsZero(),
		act.Name,
		act.Id,
	)
	if err != nil {
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rows == 0 {
		return repostore.ErrRecordNotFound
	}

	return nil
}

func (a *ActorRepository) GetActorsWithFilms(limit, offset int64) (map[*actor.Actor][]*film.Film, error) {
	result := map[*actor.Actor][]*film.Film{}
	rows, err := a.st.db.Query(
		"SELECT id, name, gender, birthdate FROM actors LIMIT $1 OFFSET $2;",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var actor actor.Actor
		err := rows.Scan(&actor.Id, &actor.Name, &actor.Gen, &actor.Birthdate)
		if err != nil {
			return nil, err
		}
		// TODO:
		films, err := a.st.FilmActor().GetActorsFilms(actor.Id)
		if err != nil {
			return nil, err
		}
		result[&actor] = films
	}

	return result, nil
}
