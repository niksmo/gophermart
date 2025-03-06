package main

import (
	"context"
	"os"
	"os/signal"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/niksmo/gophermart/internal/config"
	"github.com/niksmo/gophermart/internal/logger"
	"github.com/niksmo/gophermart/pkg/server"
	"github.com/niksmo/gophermart/pkg/sqldb"
)

func main() {
	stopCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	appLogger := logger.New()
	config.Init(appLogger)

	addrConfig := config.NewAddressConfig(appLogger)
	_ = config.NewAccrualAddrConfig(appLogger)
	dbConfig := config.NewDatabaseConfig(appLogger)
	loggerConfig := config.NewLoggerConfig(appLogger)
	logger.SetLevel(loggerConfig.Level)

	appDB := sqldb.New("pgx", dbConfig.URI, appLogger)

	appServer := server.NewHTTPServer(addrConfig.Addr(), appLogger)
	go appServer.Run()

	<-stopCtx.Done()
	appLogger.Info().Str("signal", "interrupt").Msg("shutting down gracefully")
	appServer.Close()
	appDB.Close()
}
