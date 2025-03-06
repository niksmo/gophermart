package auth

import "errors"

var (
	ErrLoginExists = errors.New("login already in use")
	ErrCredential  = errors.New("invalid login or password")
)
