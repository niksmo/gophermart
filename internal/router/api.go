package router

import (
	"context"

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

func SetupApiRoutes(ctx context.Context, appServer server.HTTPServer) {
	logging := fiberzerolog.New(fiberzerolog.Config{Logger: &logger.Instance})

	api := appServer.Group("/api", logging, compress.New())

	userPath := api.Group("/user")

	usersRepository := users.NewRepository(database.DB)
	loyaltyRepository := loyalty.NewRepository(database.DB)
	ordersRepository := orders.NewRepository(database.DB)

	// Auth
	authService := auth.NewService(
		config.Auth, usersRepository, loyaltyRepository,
	)
	authHandler := auth.NewHandler(authService)
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
	orderService := orders.NewService(ctx, ordersRepository)
	ordersHandler := orders.NewHandler(orderService)
	protectedUserPath.Post(
		"/orders",
		ordersHandler.UploadOrder,
	)
	protectedUserPath.Get(
		"/orders",
		ordersHandler.GetOrders,
	)

	// Loyalty
	loyaltyService := loyalty.NewService(ctx, loyaltyRepository)
	loyaltyHandler := loyalty.NewHandler(loyaltyService)
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
