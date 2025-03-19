package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/router"
	"github.com/niksmo/gophermart/migrations"
	"github.com/niksmo/gophermart/pkg/database"
	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/niksmo/gophermart/pkg/server"
)

func main() {
	stopCtx, stopFn := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stopFn()

	logger.Init()
	config.Init()
	loggerConfig := config.NewLoggerConfig()
	dbConfig := config.NewDatabaseConfig()
	serverConfig := config.NewServerConfig()
	authConfig := config.NewAuthConfig()
	accrualConfig := config.NewAccrualConfig()

	logger.SetLevel(loggerConfig.Level())
	database.Connect(dbConfig.URI(), logger.Instance)
	database.Migrate(migrations.Init, logger.Instance)

	appServer := server.NewHTTPServer(serverConfig.Addr(), logger.Instance)
	router.SetupAPIRoutes(stopCtx, appServer, authConfig, accrualConfig)

	go appServer.Run()

	<-stopCtx.Done()
	logger.Instance.Info().
		Str("signal", "interrupt").
		Msg("shutting down gracefully")
	appServer.Close()
	database.Close()
}
