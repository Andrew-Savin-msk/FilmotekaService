package pgstore

import (
	"context"
	"database/sql"
	"errors"

	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
)

type Users struct {
	conn *sql.DB
}

func New(dbConnectionPath string) (*Users, error) {
	conn, err := sql.Open("postgres", dbConnectionPath)
	if err != nil {
		return nil, err
	}

	err = conn.Ping()
	if err != nil {
		return nil, err
	}

	return &Users{
		conn: conn,
	}, nil
}

// TODO:
func (us Users) Create(ctx context.Context, u *user.User) error {
	stmt, err := us.conn.Prepare(
		"INSERT INTO users (email, encrypted_password) VALUES ($1, $2);",
	)

	if err != nil {
		return err
	}

	_, err = stmt.ExecContext(ctx, u.Email, u.EncPasswd)
	if err != nil {
		if errors.Is(err, sql.ErrConnDone) {
			return err
		}
		return store.ErrRecordExists
	}
	return nil
}

// TODO:
func (us *Users) Find(ctx context.Context, email string) (*user.User, error) {
	stmt, err := us.conn.Prepare(
		"SELECT id, email, ",
	)

	_, _ = stmt, err
	return nil, nil
}
