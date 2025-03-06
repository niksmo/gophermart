package auth

import (
	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
)

type AuthHandler struct {
	logger zerolog.Logger
}

func NewHandler(logger zerolog.Logger) AuthHandler {
	return AuthHandler{logger: logger}
}

func (handler AuthHandler) Register(c *fiber.Ctx) error {
	var payload SignupReqPayload
	c.BodyParser(&payload)
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(payload)
}

func (handler AuthHandler) Login(c *fiber.Ctx) error {
	var payload SigninReqPayload
	c.BodyParser(&payload)
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(payload)
}
