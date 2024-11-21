package service

import "errors"

var (
	ErrUserExists   = errors.New("user already exists")
	ErrNoSuchEntity = errors.New("entity not found")
	ErrInvalidParam = errors.New("invalid sorting parameter")
)
