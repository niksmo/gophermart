package sqldb

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type DBService struct {
	*sql.DB
	Logger zerolog.Logger
}

func New(driver, dsn string, logger zerolog.Logger) DBService {
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
	logger.Info().Str("driver", driver).Msg("database connected")
	return DBService{DB: db, Logger: logger}
}

func (s DBService) Close() {
	if err := s.DB.Close(); err != nil {
		s.Logger.Warn().Err(err).Msg("closing database connections")
	} else {
		s.Logger.Info().Msg("database connections safely closed")
	}
}
