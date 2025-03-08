package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/niksmo/gophermart/pkg/jwt"
)

const bearerPrefix = "Bearer "

type ctxKeyType struct{}

var KeyUserID ctxKeyType

type UserID int64

func Authorized(key string) fiber.Handler {
	middleware := func(c *fiber.Ctx) error {
		authorizationHeader := c.Get(fiber.HeaderAuthorization)
		if authorizationHeader == "" {
			return fiber.ErrUnauthorized
		}
		if !strings.HasPrefix(authorizationHeader, bearerPrefix) {
			return fiber.ErrUnauthorized
		}

		tokenString := strings.TrimPrefix(authorizationHeader, bearerPrefix)
		userID, err := jwt.Parse(tokenString, []byte(key))
		if err != nil {
			return fiber.ErrUnauthorized
		}
		c.Locals(KeyUserID, UserID(userID))
		return c.Next()
	}

	return middleware
}
