package database

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
)

var DB *pgxpool.Pool

func Connect(DSN string, logger zerolog.Logger) {
	var err error
	DB, err = pgxpool.New(context.Background(), DSN)
	if err != nil {
		logger.Fatal().Err(err).Caller().Msg("unable to create connection pool")
	}
}

func Migrate(stmt string, logger zerolog.Logger) {
	tag, err := DB.Exec(context.Background(), stmt)
	if err != nil {
		logger.Fatal().Err(err).Caller()
		return
	}
	logger.Info().Str("tag", tag.String()).Msg("database migration")
}

func Close() {
	DB.Close()
}
