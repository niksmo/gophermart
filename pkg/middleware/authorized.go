package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/pkg/jwt"
)

const bearerPrefix = "Bearer "

type ctxKeyType struct{}

var KeyUserID ctxKeyType

type UserID int64

func (userID UserID) Int32() int32 {
	return int32(userID)
}

func Authorized(key []byte) fiber.Handler {
	middleware := func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(fiber.HeaderAuthorization)
		if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
			return fiber.ErrUnauthorized
		}

		tokenString := strings.TrimPrefix(authorizationHeader, bearerPrefix)
		userID, err := jwt.Parse(tokenString, key)
		if err != nil {
			return fiber.ErrUnauthorized
		}
		c.Locals(KeyUserID, UserID(userID))
		return c.Next()
	}

	return middleware
}

func GetUserID(c *fiber.Ctx) (UserID, error) {
	userID, ok := c.Locals(KeyUserID).(UserID)
	if !ok {
		return userID, errors.New("extracting userID from fiber.Ctx.Locals")
	}
	return userID, nil
}
