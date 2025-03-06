package repository

import "errors"

var (
	ErrExists   = errors.New("entry already exists")
	ErrNotFound = errors.New("entry not found")
)
