package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/internal/auth"
	"github.com/niksmo/gophermart/internal/middleware"
)

func SetUserPath(router fiber.Router) {
	authHandler := auth.NewHandler()
	router = router.Group("/user", middleware.AllowJSON)
	router.Post("/register", authHandler.Register)
	router.Post("/login", authHandler.Login)
}
