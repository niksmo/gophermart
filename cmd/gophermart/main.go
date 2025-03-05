package main

import (
	"log"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/niksmo/gophermart/internal/config"
	"github.com/niksmo/gophermart/internal/logger"
	"github.com/niksmo/gophermart/pkg/sqldb"
)

func main() {
	appLogger := logger.New()
	config.Init(appLogger)

	addrConfig := config.NewAddressConfig()
	accrualConfig := config.NewAccrualAddrConfig()
	dbConfig := config.NewDatabaseConfig()
	loggerConfig := config.NewLoggerConfig(appLogger)
	logger.SetLevel(loggerConfig.Level)

	log.Println("addr:", addrConfig)
	log.Println("accrual addr:", accrualConfig.GetOrdersReqURL("54321"))
	log.Println("db uri:", dbConfig.URI)

	pgDB := sqldb.New("pgx", dbConfig.URI)
	defer pgDB.Close()
}
