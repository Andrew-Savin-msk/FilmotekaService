package apiserver

import "errors"

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not auntificated")
	errResourceForbiden         = errors.New("you dont have permossions to get this resource")
	errIncorrectId              = errors.New("presented incorrect id type")
	errMethodNotAllowed         = errors.New("unsuportable method type")

	ErrInternalDbError     = errors.New("valid ending of operation is unable")
	ErrInvalidCredentials  = errors.New("invalid login")
	ErrHashingPassword     = errors.New("unable to hash password")
	ErrPanicHanding        = errors.New("internal server troubles")
	ErrMissingToken        = errors.New("missing authorization token")
	ErrExpiredToken        = errors.New("token epired")
	ErrNotAuntificated     = errors.New("user is not auntificated")
	ErrInvalidToken        = errors.New("invalid acess token")
	ErrUnsupportedMethod   = errors.New("this method doesn't supported")
	ErrInvalidQuerryParams = errors.New("invalid or not existing querry parameters")
	ErrInvalidRequestBody  = errors.New("invalid request body")
	ErrServiceUnavailable  = errors.New("service currently is not available")
	ErrNoSuchResorce       = errors.New("resorce doesn't exist")
)
