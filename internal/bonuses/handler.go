package bonuses

import (
	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/niksmo/gophermart/pkg/middleware"
)

type BonusesHandler struct {
	service BonusesService
}

func NewBonusesHandler(service BonusesService) BonusesHandler {
	return BonusesHandler{service: service}
}

func (h BonusesHandler) ShowBalance(c *fiber.Ctx) error {
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
