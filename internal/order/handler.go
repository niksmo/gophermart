package order

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
)

type OrderHandler struct {
	service OrderService
}

func NewHandler(service OrderService) OrderHandler {
	return OrderHandler{service: service}
}

func (handler OrderHandler) UploadOrder(c *fiber.Ctx) error {
	payload := UploadOrderReqPayload(string(c.Body()))
	number, err := payload.Validate()
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.SendString(strconv.FormatInt(number, 10))
}
