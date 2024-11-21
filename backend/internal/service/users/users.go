package usersservice

import (
	"context"
	"errors"

	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/service"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/store"
	"github.com/sirupsen/logrus"
)

type Users struct {
	users store.Users

	logger *logrus.Entry
	ctx    context.Context
}

func New(ctx context.Context, logger *logrus.Entry, users store.Users) service.Users {
	return &Users{
		users: users,
		logger: logger.WithFields(logrus.Fields{
			"layer":     "service",
			"structure": "users",
		}),
		ctx: ctx,
	}
}

func (us *Users) Create(u *user.User) error {
	err := us.users.Create(us.ctx, u)
	if err != nil {
		if errors.Is(err, store.ErrRecordExists) {
			return service.ErrUserExists
		}
		us.logger.Errorf("op: Create, unexpected error: %s", err)
		return err
	}
	return nil
}

func (us *Users) GetByEmail(email string) (*user.User, error) {
	u, err := us.users.Find(us.ctx, email)
	if err != nil {
		if errors.Is(err, store.ErrRecordNotFound) {
			return nil, service.ErrNoSuchEntity
		}
		us.logger.Errorf("op: GetByEmail, unexpected error: %s", err)
		return nil, err
	}
	return u, nil
}
