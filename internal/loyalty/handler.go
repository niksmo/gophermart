package loyalty

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/niksmo/gophermart/pkg/middleware"
)

type LoyaltyHandler struct {
	service LoyaltyService
}

func NewHandler(service LoyaltyService) LoyaltyHandler {
	return LoyaltyHandler{service: service}
}

func (h LoyaltyHandler) GetBalance(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}
	balance, err := h.service.GetUserBalance(c.Context(), userID.Int32())
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}
	return c.JSON(balance)
}

func (h LoyaltyHandler) WithdrawPoints(c *fiber.Ctx) error {
	var payload WithdrawRequestScheme
	err := c.BodyParser(&payload)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	validationData, ok := payload.Validate()
	if !ok {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(validationData)
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}

	err = h.service.WithdrawPoints(
		c.Context(), userID.Int32(), payload.OrderNumber, payload.Amount,
	)

	if err != nil {
		if errors.Is(err, errs.ErrLoyaltyNotEnoughPoints) {
			return fiber.NewError(fiber.StatusPaymentRequired, err.Error())
		}
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}
	return c.SendStatus(fiber.StatusOK)
}

func (h LoyaltyHandler) GetWithdrawals(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}

	withdrawals, err := h.service.GetUserWithdrawals(c.Context(), userID.Int32())
	if err != nil {
		if errors.Is(err, errs.ErrLoyaltyNoWithdrawals) {
			return fiber.NewError(fiber.StatusNoContent, err.Error())
		}
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}

	return c.JSON(withdrawals)
}
