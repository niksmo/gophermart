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

func SetupAPIRoutes(
	ctx context.Context,
	appServer server.HTTPServer,
	authConfig config.AuthConfig,
	accrualConfig config.AccrualConfig,
) {
	logMdw := fiberzerolog.New(
		fiberzerolog.Config{Logger: &logger.Instance},
	)
	compressMdw := compress.New()

	api := appServer.Group("/api", logMdw, compressMdw)

	userPath := api.Group("/user")

	usersRepository := users.NewRepository(database.DB)
	loyaltyRepository := loyalty.NewRepository(database.DB)
	ordersRepository := orders.NewRepository(database.DB)

	// Auth
	authService := auth.NewService(
		authConfig, usersRepository, loyaltyRepository,
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
		"", middleware.Authorized(authConfig.Key()),
	)

	// Orders
	ordersService := orders.NewService(ctx, ordersRepository, accrualConfig)
	go ordersService.Restore(ctx)
	go ordersService.FlushAccrualResults(ctx)

	ordersHandler := orders.NewHandler(ordersService)
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
