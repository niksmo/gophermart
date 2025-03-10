package router

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/auth"
	"github.com/niksmo/gophermart/internal/bonuses"
	"github.com/niksmo/gophermart/internal/orders"
	"github.com/niksmo/gophermart/internal/users"
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
		auth.NewService(
			config.Auth,
			users.NewRepository(database.DB),
			bonuses.NewRepository(database.DB),
		),
	)
	api.Post("/user/register", middleware.RequireJSON, authHandler.Register)
	api.Post("/user/login", middleware.RequireJSON, authHandler.Login)

	authorized := api.Group("", middleware.Authorized(config.Auth.Key()))

	// Orders
	orderHandler := orders.NewHandler(
		orders.NewService(orders.NewRepository(database.DB)),
	)
	authorized.Post("/user/orders", orderHandler.UploadOrder)
	authorized.Get("/user/orders", orderHandler.GetOrders)
}
