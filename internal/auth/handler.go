package auth

import (
	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
}

func NewHandler() AuthHandler {
	return AuthHandler{}
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
