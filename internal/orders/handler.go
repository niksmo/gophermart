package orders

import (
	"errors"

	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/internal/errs"
	"github.com/niksmo/gophermart/pkg/logger"
	"github.com/niksmo/gophermart/pkg/middleware"
)

type OrdersHandler struct {
	service OrdersService
}

func NewHandler(service OrdersService) OrdersHandler {
	return OrdersHandler{service: service}
}

func (h OrdersHandler) UploadOrder(c *fiber.Ctx) error {
	payload := OrderNumberScheme(c.Body())
	orderNumber, err := payload.Validate()
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderInvalidFormat):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case errors.Is(err, errs.ErrOrderInvalidNum):
			return fiber.NewError(
				fiber.StatusUnprocessableEntity, err.Error(),
			)
		default:
			logger.Instance.Error().Err(err).Caller().Send()
			return fiber.ErrInternalServerError
		}
	}

	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}

	err = h.service.UploadOrder(c.Context(), userID.Int32(), orderNumber)
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderUploadedByUser):
			return fiber.NewError(fiber.StatusOK, err.Error())
		case errors.Is(err, errs.ErrOrderUploadedByOther):
			return fiber.NewError(fiber.StatusConflict, err.Error())
		default:
			logger.Instance.Error().Err(err).Caller().Send()
			return fiber.ErrInternalServerError
		}
	}
	return c.SendStatus(fiber.StatusAccepted)
}

func (h OrdersHandler) GetOrders(c *fiber.Ctx) error {
	userID, err := middleware.GetUserID(c)
	if err != nil {
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}

	orderList, err := h.service.GetUserOrders(c.Context(), userID.Int32())
	if err != nil {
		if errors.Is(err, errs.ErrOrdersNoUploads) {
			return fiber.NewError(fiber.StatusNoContent, err.Error())
		}
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}
	return c.JSON(orderList)
}
