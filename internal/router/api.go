package router

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/auth"
	"github.com/niksmo/gophermart/internal/order"
	"github.com/niksmo/gophermart/internal/repository"
	"github.com/niksmo/gophermart/pkg/database"
	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/niksmo/gophermart/pkg/middleware"
	"github.com/niksmo/gophermart/pkg/server"
)

func SetupApiRoutes(appServer server.HTTPServer) {
	logging := fiberzerolog.New(fiberzerolog.Config{Logger: &logger.Instance})

	api := appServer.Group("/api", logging, compress.New())

	// Auth
	authHandler := auth.NewHandler(
		auth.NewService(config.Auth, repository.Users(database.DB)),
	)
	api.Post("/user/register", middleware.RequireJSON, authHandler.Register)
	api.Post("/user/login", middleware.RequireJSON, authHandler.Login)

	requireAuth := middleware.Authorized(config.Auth.Key())

	// Order
	orderHandler := order.NewHandler(
		order.NewService(repository.Orders(database.DB)),
	)
	api.Post("/user/orders", requireAuth, orderHandler.UploadOrder)
}
