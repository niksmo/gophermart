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
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(NewSignupResPayload(tokenString))
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	var payload SigninReqPayload
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
	c.Set(fiber.HeaderCacheControl, "no-store")
	return c.JSON(NewSigninResPayload(tokenString))
}
