package auth

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/di"
	"github.com/niksmo/gophermart/pkg/logger"
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
	if err := validateAuthRequest(c, &payload); err != nil {
		return err
	}

	tokenString, err := h.service.RegisterUser(
		c.Context(), payload.Login, payload.Password,
	)
	if err != nil {
		if errors.Is(err, errs.ErrUserLoginExists) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}
	return authorize(c, tokenString)
}

func (h AuthHandler) Login(c *fiber.Ctx) error {
	var payload SigninRequestScheme
	if err := validateAuthRequest(c, &payload); err != nil {
		return err
	}

	tokenString, err := h.service.AuthorizeUser(
		c.Context(), payload.Login, payload.Password,
	)

	if err != nil {
		if errors.Is(err, errs.ErrUserCredentials) {
			return fiber.NewError(fiber.StatusUnauthorized, err.Error())
		}
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}
	return authorize(c, tokenString)
}

func validateAuthRequest(c *fiber.Ctx, payload any) error {
	err := c.BodyParser(payload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	validator, ok := payload.(di.Validator)
	if !ok {
		return fiber.ErrInternalServerError
	}

	result, valid := validator.Validate()
	if !valid {
		return c.Status(fiber.StatusBadRequest).JSON(result)
	}
	return nil
}

func authorize(c *fiber.Ctx, tokenString string) error {
	c.Set(fiber.HeaderCacheControl, "no-store")
	c.Set(fiber.HeaderAuthorization, tokenType+tokenString)
	return c.SendStatus(fiber.StatusOK)
}
