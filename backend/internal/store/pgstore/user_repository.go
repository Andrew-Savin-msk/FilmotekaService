package pgstore

import (
	"database/sql"

	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
)

type UserRepository struct {
	st *Store
}

func (r *UserRepository) Create(u *user.User) error {
	// Encrypt
	err := u.Sequre()
	if err != nil {
		return err
	}

	// Save
	return r.st.db.QueryRow(
		"INSERT INTO users (email, encrypted_password, is_admin) VALUES ($1, $2, $3) RETURNING id",
		u.Email,
		u.EncPasswd,
		u.IsAdmin,
	).Scan(&u.Id)
}

func (r *UserRepository) Find(id int) (*user.User, error) {
	u := &user.User{
		Id: id,
	}
	err := r.st.db.QueryRow(
		"SELECT email, encrypted_password, is_admin FROM users WHERE id = $1",
		id,
	).Scan(
		&u.Email,
		&u.EncPasswd,
		&u.IsAdmin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}

func (r *UserRepository) FindByEmail(email string) (*user.User, error) {
	u := &user.User{
		Email: email,
	}
	err := r.st.db.QueryRow(
		"SELECT id, encrypted_password, is_admin FROM users WHERE email = $1",
		email,
	).Scan(
		&u.Id,
		&u.EncPasswd,
		&u.IsAdmin,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		return nil, err
	}
	return u, nil
}
