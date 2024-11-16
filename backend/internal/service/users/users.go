package userservice

import (
	"context"

	user "github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/model/user"
	"github.com/Andrew-Savin-msk/filmoteka-service/backend/internal/service"
	"github.com/sirupsen/logrus"
)

// TODO:
type Users struct {
	// Users repo/table connection

	logger *logrus.Entry
	ctx    context.Context
}

func New(ctx context.Context, logger *logrus.Entry) service.Users {
	return &Users{
		logger: logger,
		ctx:    ctx,
	}
}

func (us *Users) Create(u *user.User) error {
	// db Create method
	panic("unimplemented")
}

func (us *Users) GetByEmail(email string) (*user.User, error) {
	// db FindByEmail method
	panic("unimplemented")
}
