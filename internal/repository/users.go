package repository

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/logger"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func Users(db *pgxpool.Pool) UsersRepository {
	return UsersRepository{db: db}
}

func (r UsersRepository) Create(
	ctx context.Context, login string, password string,
) error {
	stmt := `
	INSERT INTO users (login, password) VALUES ($1, $2);
	`
	_, err := r.db.Exec(ctx, stmt, login, password)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == pgerrcode.UniqueViolation {
				return errs.ErrLoginExists
			}
			logger.Instance.Warn().Err(err).Msg("creating user")
			return err
		}
	}
	logger.Instance.Info().Msg("user created")
	return nil
}

func (r UsersRepository) ReadByLogin(ctx context.Context, login string) {}

func (r UsersRepository) ReadByID(ctx context.Context, userID int64) {}
