package router

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/auth"
	"github.com/niksmo/gophermart/internal/repository"
	"github.com/niksmo/gophermart/pkg/database"
	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/niksmo/gophermart/pkg/middleware"
	"github.com/niksmo/gophermart/pkg/server"
)

func SetupApiRoutes(appServer server.HTTPServer) {
	api := appServer.Group(
		"/api",
		fiberzerolog.New(fiberzerolog.Config{Logger: &logger.Instance}),
		middleware.AllowJSON,
	)

	_ = middleware.Authorized(config.Auth.Key())

	// Auth
	authHandler := auth.NewHandler(
		auth.NewService(config.Auth, repository.Users(database.DB)),
	)
	api.Post("/user/register", authHandler.Register)
	api.Post("/user/login", authHandler.Login)
}
