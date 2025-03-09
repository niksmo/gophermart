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
	stopCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	logger.Init()
	config.Init()
	logger.SetLevel(config.Logger.Level())
	database.Connect(config.Database.URI(), logger.Instance)
	database.Migrate(migrations.Init, logger.Instance)

	appServer := server.NewHTTPServer(config.Server.Addr(), logger.Instance)
	router.SetupApiRoutes(appServer)

	go appServer.Run()
	<-stopCtx.Done()
	logger.Instance.Info().Str("signal", "interrupt").Msg("shutting down gracefully")
	appServer.Close()
	database.Close()
}
