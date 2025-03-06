package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/internal/auth"
	"github.com/niksmo/gophermart/internal/middleware"
	"github.com/rs/zerolog"
)

func SetUserPath(router fiber.Router, logger zerolog.Logger) {
	authHandler := auth.NewHandler(logger)
	router = router.Group("/user", middleware.AllowJSON)
	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)
}
