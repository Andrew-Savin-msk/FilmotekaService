package store

import "errors"

var (
	ErrRecordExists   = errors.New("record already exists")
	ErrRecordNotFound = errors.New("record not found")
	ErrInvalidParam   = errors.New("invalid query parameter")
)
