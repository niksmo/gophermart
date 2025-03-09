package router

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/auth"
	"github.com/niksmo/gophermart/internal/repository"
	"github.com/niksmo/gophermart/pkg/database"
	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/niksmo/gophermart/pkg/middleware"
	"github.com/niksmo/gophermart/pkg/server"
)

func SetupApiRoutes(appServer server.HTTPServer) {
	logging := fiberzerolog.New(fiberzerolog.Config{Logger: &logger.Instance})

	api := appServer.Group(
		"/api",
		logging,
		middleware.AllowJSON,
		compress.New(),
	)

	// Auth
	authHandler := auth.NewHandler(
		auth.NewService(config.Auth, repository.Users(database.DB)),
	)
	api.Post("/user/register", authHandler.Register)
	api.Post("/user/login", authHandler.Login)

	_ = middleware.Authorized(config.Auth.Key())
}
