package errs

import "errors"

var (
	ErrLoginExists = errors.New("user login already in use")
	ErrCredentials = errors.New("invalid login or password")
)
