package router

import (
	"context"

	"github.com/gofiber/contrib/fiberzerolog"
	"github.com/gofiber/fiber/v2"
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

func SetupAPIRoute(
	ctx context.Context,
	appServer server.HTTPServer,
	authConfig config.AuthConfig,
	accrualConfig config.AccrualConfig,
) {
	apiPath := setAPIPath(appServer)

	userPath := apiPath.Group("/user")

	usersRepository := users.NewRepository(database.DB)
	loyaltyRepository := loyalty.NewRepository(database.DB)
	ordersRepository := orders.NewRepository(database.DB)

	setAuthPaths(userPath, authConfig, usersRepository, loyaltyRepository)

	protectedRouter := userPath.Group(
		"", middleware.Authorized(authConfig.Key()),
	)

	setOrdersPaths(ctx, protectedRouter, accrualConfig, ordersRepository)
	setLoyaltyPaths(protectedRouter, loyaltyRepository)
}

func setAPIPath(appServer server.HTTPServer) fiber.Router {
	logMdw := fiberzerolog.New(
		fiberzerolog.Config{Logger: &logger.Instance},
	)
	compressMdw := compress.New()

	return appServer.Group("/api", logMdw, compressMdw)
}

func setAuthPaths(
	router fiber.Router,
	authConfig config.AuthConfig,
	usersRepository users.UsersRepository,
	loyaltyRepository loyalty.LoyaltyRepository,
) {

	authService := auth.NewService(
		authConfig, usersRepository, loyaltyRepository,
	)
	authHandler := auth.NewHandler(authService)
	router.Post(
		"/register",
		middleware.RequireJSON,
		authHandler.Register,
	)
	router.Post(
		"/login",
		middleware.RequireJSON,
		authHandler.Login,
	)
}

func setOrdersPaths(
	ctx context.Context,
	router fiber.Router,
	accrualConfig config.AccrualConfig,
	ordersRepository orders.OrdersRepository,
) {
	ordersService := orders.NewService(ctx, ordersRepository, accrualConfig)
	go ordersService.Restore(ctx)
	go ordersService.FlushAccrualResults(ctx)

	ordersHandler := orders.NewHandler(ordersService)
	router.Post(
		"/orders",
		ordersHandler.UploadOrder,
	)
	router.Get(
		"/orders",
		ordersHandler.GetOrders,
	)
}

func setLoyaltyPaths(
	router fiber.Router, loyaltyRepository loyalty.LoyaltyRepository,
) {
	loyaltyService := loyalty.NewService(loyaltyRepository)
	loyaltyHandler := loyalty.NewHandler(loyaltyService)
	router.Get(
		"/balance",
		loyaltyHandler.GetBalance,
	)

	router.Post(
		"/balance/withdraw",
		middleware.RequireJSON,
		loyaltyHandler.WithdrawPoints,
	)

	router.Get(
		"/withdrawals",
		loyaltyHandler.GetWithdrawals,
	)
}
