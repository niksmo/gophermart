package errs

import "errors"

var (
	ErrUserLoginLength = errors.New("length min 1 and max 60")

	ErrUserLoginInvalid = errors.New(
		"allow latin charecters, numbers, dash, underline",
	)

	ErrUserPasswordLength = errors.New("length min 1 and max 60")

	ErrUserPasswordInvalid = errors.New(
		"allow latin charecters, special symbols, without spaces",
	)

	ErrUserLoginExists = errors.New("user login already in use")

	ErrUserCredentials = errors.New("invalid login or password")
)
