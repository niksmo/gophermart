package loyalty

import (
	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/niksmo/gophermart/pkg/middleware"
)

type LoyaltyHandler struct {
	service LoyaltyService
}

func NewHandler(service LoyaltyService) LoyaltyHandler {
	return LoyaltyHandler{service: service}
}

func (h LoyaltyHandler) ShowBalance(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}
	balance, err := h.service.GetUsersBalance(c.Context(), userID.Int32())
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}
	return c.JSON(balance)
}
