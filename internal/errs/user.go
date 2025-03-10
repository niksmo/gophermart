package errs

import "errors"

var (
	ErrLoginLength  = errors.New("length min 1 and max 60")
	ErrLoginInvalid = errors.New("allow latin charecters, numbers, dash, underlining")

	ErrPasswordLength  = errors.New("length min 1 and max 60")
	ErrPasswordInvalid = errors.New("allow latin charecters, special symbols, without spaces")

	ErrLoginExists = errors.New("user login already in use")
	ErrCredentials = errors.New("invalid login or password")
)
