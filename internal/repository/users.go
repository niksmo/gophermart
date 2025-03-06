package repository

import (
	"context"
	"database/sql"

	"github.com/rs/zerolog"
)

type UsersRepository struct {
	logger zerolog.Logger
	db     *sql.DB
}

func Users(db *sql.DB, logger zerolog.Logger) UsersRepository {
	return UsersRepository{db: db, logger: logger}
}

func (repo UsersRepository) Create(
	ctx context.Context, login string, password string,
) {
}

func (repo UsersRepository) ReadByLogin(ctx context.Context, login string) {}

func (repo UsersRepository) ReadByID(ctx context.Context, userID int64) {}
