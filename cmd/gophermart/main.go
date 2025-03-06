package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/gofiber/contrib/fiberzerolog"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/niksmo/gophermart/internal/api"
	"github.com/niksmo/gophermart/internal/config"
	"github.com/niksmo/gophermart/internal/logger"
	"github.com/niksmo/gophermart/internal/repository"
	"github.com/niksmo/gophermart/pkg/server"
	"github.com/niksmo/gophermart/pkg/sqldb"
)

func main() {
	stopCtx, _ := signal.NotifyContext(context.Background(), os.Interrupt)
	logger.Init()
	config.Init()

	addrConfig := config.NewAddressConfig()
	_ = config.NewAccrualAddrConfig()
	dbConfig := config.NewDatabaseConfig()
	loggerConfig := config.NewLoggerConfig()
	logger.SetLevel(loggerConfig.Level)

	appDB := sqldb.New("pgx", dbConfig.URI, logger.Instance)
	repository.Init(appDB)

	appServer := server.NewHTTPServer(addrConfig.Addr(), logger.Instance)
	appServer.Use(fiberzerolog.New(fiberzerolog.Config{Logger: &logger.Instance}))

	router := appServer.Group("/api")
	api.SetUserPath(router, appDB)

	go appServer.Run()

	<-stopCtx.Done()
	logger.Instance.Info().Str("signal", "interrupt").Msg("shutting down gracefully")
	appServer.Close()
	appDB.Close()
}
