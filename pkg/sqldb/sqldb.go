package sqldb

import (
	"database/sql"

	"github.com/rs/zerolog"
)

type SQLDB struct {
	*sql.DB
	logger zerolog.Logger
}

func New(driver, dsn string, logger zerolog.Logger) SQLDB {
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
	return SQLDB{DB: db, logger: logger}
}

func (sqldb SQLDB) Close() {
	if err := sqldb.DB.Close(); err != nil {
		sqldb.logger.Warn().Err(err).Msg("closing database connections")
	} else {
		sqldb.logger.Info().Msg("database connections safely closed")
	}
}
