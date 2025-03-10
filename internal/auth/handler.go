package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/internal/errs"
)

const tokenType = "Bearer "

type AuthHandler struct {
	service AuthService
}

func NewHandler(service AuthService) AuthHandler {
	return AuthHandler{service: service}
}

func (h AuthHandler) Register(c *fiber.Ctx) error {
	var payload SignupRequestScheme
	err := c.BodyParser(&payload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	validationResult, ok := payload.Validate()
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(validationResult)
	}

	tokenString, err := h.service.RegisterUser(
		c.Context(), payload.Login, payload.Password,
	)
	if err != nil {
		if errors.Is(err, errs.ErrLoginExists) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.ErrInternalServerError
	}
	return authorize(c, tokenString)
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	var payload SigninRequestScheme
	err := c.BodyParser(&payload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	validationResult, ok := payload.Validate()
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(validationResult)
	}

	tokenString, err := h.service.AuthorizeUser(
		c.Context(), payload.Login, payload.Password,
	)

	if err != nil {
		if errors.Is(err, errs.ErrCredentials) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		return fiber.ErrInternalServerError
	}
	return authorize(c, tokenString)
}

func authorize(c *fiber.Ctx, tokenString string) error {
	c.Set(fiber.HeaderCacheControl, "no-store")
	c.Set(fiber.HeaderAuthorization, tokenType+tokenString)
	return c.SendStatus(fiber.StatusOK)
}
