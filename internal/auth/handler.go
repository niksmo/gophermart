package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/internal/errs"
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
	err := h.service.RegisterUser(c.Context(), payload.Login, payload.Password)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrLoginExists):
			return fiber.NewError(fiber.StatusConflict, errs.ErrLoginExists.Error())
		default:
			return fiber.ErrInternalServerError
		}
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.SendString("Registered")
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	var payload SigninReqPayload
	c.BodyParser(&payload)
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(payload)
}
