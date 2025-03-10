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

func (handler OrdersHandler) UploadOrder(c *fiber.Ctx) error {
	payload := OrderNumberScheme(c.Body())
	orderNumber, err := payload.Validate()
	if err != nil {
		switch {
		case errors.Is(err, errs.ErrOrderInvalidFormat):
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		case errors.Is(err, errs.ErrOrderInvalidNum):
			return fiber.NewError(fiber.StatusUnprocessableEntity, err.Error())
		default:
			logger.Instance.Error().Err(err).Caller().Send()
			return fiber.ErrInternalServerError
		}
	}

	userID, ok := middleware.GetUserID(c)
	if !ok {
		logger.Instance.Error().Err(err).Caller().Send()
		return fiber.ErrInternalServerError
	}

	err = handler.service.UploadOrder(
		c.Context(), userID.Int32(), orderNumber,
	)
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
