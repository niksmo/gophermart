package router

import (
	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2/middleware/compress"
	"github.com/niksmo/gophermart/config"
	"github.com/niksmo/gophermart/internal/auth"
	"github.com/niksmo/gophermart/internal/loyalty"
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

	userPath := api.Group("/user")

	// Auth
	authHandler := auth.NewHandler(
		auth.NewService(
			config.Auth,
			users.NewRepository(database.DB),
			loyalty.NewRepository(database.DB),
		),
	)
	userPath.Post(
		"/register",
		middleware.RequireJSON,
		authHandler.Register,
	)
	userPath.Post(
		"/login",
		middleware.RequireJSON,
		authHandler.Login,
	)

	protectedUserPath := userPath.Group(
		"", middleware.Authorized(config.Auth.Key()),
	)

	// Orders
	ordersHandler := orders.NewHandler(
		orders.NewService(orders.NewRepository(database.DB)),
	)
	protectedUserPath.Post(
		"/orders",
		ordersHandler.UploadOrder,
	)
	protectedUserPath.Get(
		"/orders",
		ordersHandler.GetOrders,
	)

	// Loyalty
	loyaltyHandler := loyalty.NewHandler(
		loyalty.NewService(loyalty.NewRepository(database.DB)),
	)
	protectedUserPath.Get(
		"/balance",
		loyaltyHandler.GetBalance,
	)

	protectedUserPath.Post(
		"/balance/withdraw",
		middleware.RequireJSON,
		loyaltyHandler.WithdrawPoints,
	)

	protectedUserPath.Get(
		"/withdrawals",
		loyaltyHandler.GetWithdrawals,
	)
}
