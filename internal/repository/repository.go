package repository

import (
	"github.com/niksmo/gophermart/internal/logger"
	"github.com/niksmo/gophermart/migrations"
	"github.com/niksmo/gophermart/pkg/sqldb"
)

func Init(dbService sqldb.DBService) {
	_, err := dbService.Exec(migrations.Init)
	if err != nil {
		logger.Instance.Fatal().Err(err).Caller().Msg("repository initializing")
	} else {
		logger.Instance.Info().Msg("repository initialized")
	}
}
