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
	var reqPayload SignupReqPayload
	err := c.BodyParser(&reqPayload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	validationResult, ok := reqPayload.Validate()
	if !ok {
		return c.Status(fiber.StatusBadRequest).JSON(validationResult)
	}

	tokenString, err := h.service.RegisterUser(
		c.Context(), reqPayload.Login, reqPayload.Password,
	)
	if err != nil {
		if errors.Is(err, errs.ErrLoginExists) {
			return fiber.NewError(fiber.StatusConflict, err.Error())
		}
		return fiber.ErrInternalServerError
	}
	resPayload := NewResPayload(tokenString)
	c.Set(fiber.HeaderCacheControl, "no-store")
	c.Set(fiber.HeaderAuthorization, resPayload.String())
	return c.JSON(resPayload)
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
	resPayload := NewResPayload(tokenString)
	c.Set(fiber.HeaderCacheControl, "no-store")
	c.Set(fiber.HeaderAuthorization, resPayload.String())
	return c.JSON(resPayload)
}
