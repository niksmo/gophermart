package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/internal/auth"
	"github.com/niksmo/gophermart/internal/config"
	"github.com/niksmo/gophermart/internal/middleware"
	"github.com/niksmo/gophermart/internal/repository"
	"github.com/niksmo/gophermart/pkg/sqldb"
)

func SetUserPath(router fiber.Router, authConfig config.AuthConfig, dbService sqldb.DBService) {
	router = router.Group("/user", middleware.AllowJSON)

	authHandler := auth.NewHandler(
		auth.NewService(authConfig, repository.Users(dbService)),
	)
	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)
}
