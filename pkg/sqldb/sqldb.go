package sqldb

import (
	"database/sql"

	"github.com/rs/zerolog"
)

func New(driver, dsn string, logger *zerolog.Logger) *sql.DB {
	logFatal := func(err error) {
		logger.Fatal().Err(err).Caller().Send()
	}

	db, err := sql.Open(driver, dsn)

	if err != nil {
		logFatal(err)
	}
	if err = db.Ping(); err != nil {
		logFatal(err)
	}
	logger.Info().Str("driver", driver).Msg("databse is connected")
	return db
}
