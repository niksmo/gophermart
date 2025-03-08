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
	err := c.BodyParser(&payload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// validate payload
	// if err return fiber.Err

	tokenString, err := h.service.RegisterUser(
		c.Context(), payload.Login, payload.Password,
	)
	if err != nil {
		if errors.Is(err, errs.ErrLoginExists) {
			return fiber.NewError(fiber.StatusConflict, errs.ErrLoginExists.Error())
		}
		return fiber.ErrInternalServerError
	}
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(NewSignupResPayload(tokenString))
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	var payload SigninReqPayload
	err := c.BodyParser(&payload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	// validate payload
	// if err return fiber.Err

	// h.service.AuthorizeUser(c.Context(), payload.Login, payload.Password)
	// process errors

	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(NewSigninResPayload("<TokenValue>"))
}
