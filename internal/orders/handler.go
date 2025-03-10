package orders

import (
	"github.com/gofiber/fiber/v2"
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
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	userID, ok := c.Locals(middleware.KeyUserID).(middleware.UserID)
	if !ok {
		return fiber.ErrInternalServerError
	}

	err = handler.service.UploadOrder(
		c.Context(), userID.Int64(), orderNumber,
	)
	if err != nil {
		// handle errors
	}
	return c.SendStatus(fiber.StatusAccepted)
}
