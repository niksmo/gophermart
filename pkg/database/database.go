package database

import (
	"context"
	"errors"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

var DB *pgxpool.Pool

func Connect(DSN string, logger zerolog.Logger) {
	var err error
	DB, err = pgxpool.New(context.Background(), DSN)
	if err != nil {
		logger.Fatal().Err(err).Caller().Msg("unable to create connection pool")
		return
	}
	err = DB.Ping(context.Background())
	if err != nil {
		logger.Fatal().Err(err).Caller().Msg("database not connected")
		return
	}
	logger.Info().Msg("database connected")
}

func Migrate(stmt string, logger zerolog.Logger) {
	tag, err := DB.Exec(context.Background(), stmt)
	if err != nil {
		logger.Fatal().Err(err).Caller().Send()
		return
	}
	logger.Info().Str("tag", tag.String()).Msg("database migration")
}

func Close() {
	DB.Close()
}

func IsUniqueError(err error) bool {
	var pgErr *pgconn.PgError
	return errors.As(err, &pgErr) &&
		pgErr.Code == pgerrcode.UniqueViolation
}
