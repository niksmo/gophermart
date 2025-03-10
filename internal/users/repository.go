package users

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/logger"
)

type UsersRepository struct {
	db *pgxpool.Pool
}

func NewRepository(db *pgxpool.Pool) UsersRepository {
	return UsersRepository{db: db}
}

func (r UsersRepository) Create(
	ctx context.Context, login string, password string,
) (int64, error) {
	stmt := `
	INSERT INTO users (login, password) VALUES ($1, $2)
	RETURNING id;
	`
	var userID int64
	err := r.db.QueryRow(ctx, stmt, login, password).Scan(&userID)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == pgerrcode.UniqueViolation {
			return -1, errs.ErrUserLoginExists
		}
		logger.Instance.Warn().Err(err).Msg("creating user")
		return -1, err
	}
	return userID, nil
}

func (r UsersRepository) ReadByLogin(
	ctx context.Context, login string,
) (int64, string, error) {
	stmt := `
	SELECT id, password FROM users WHERE login=$1;
	`
	var (
		userID  int64
		pwdHash string
	)
	err := r.db.QueryRow(ctx, stmt, login).Scan(&userID, &pwdHash)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return -1, "", errs.ErrUserCredentials
		}
		logger.Instance.Warn().Err(err).Msg("reading by login")
		return -1, "", err
	}
	return userID, pwdHash, nil
}
