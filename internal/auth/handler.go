package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

type AuthHandler struct {
	service AuthService
}

func NewHandler(service AuthService) AuthHandler {
	return AuthHandler{service: service}
}

func (h AuthHandler) Register(c *fiber.Ctx) error {
	var payload SignupReqPayload
	c.BodyParser(&payload)
	c.Set(fiber.HeaderCacheControl, "no-store")
	err := h.service.RegisterUser(c.Context(), payload.Login, payload.Password)
	if errors.Is(err, ErrLoginExists) {
		return fiber.NewError(fiber.StatusConflict, ErrLoginExists.Error())
	}
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.SendString("Registered")
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	var payload SigninReqPayload
	c.BodyParser(&payload)
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(payload)
}
