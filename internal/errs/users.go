package errs

import "errors"

var (
	ErrLoginExists = errors.New("user login already in use")
	ErrCredential  = errors.New("invalid login or password")
)
