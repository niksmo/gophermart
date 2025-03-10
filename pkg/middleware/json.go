package middleware

import "github.com/gofiber/fiber/v2"

func RequireJSON(c *fiber.Ctx) error {
	if unsafeMethod(c.Method()) {
		if c.Get(fiber.HeaderContentType) != fiber.MIMEApplicationJSON {
			return fiber.ErrUnsupportedMediaType
		}
		return c.Next()
	}
	return c.Next()
}

var unsafeMethods = map[string]struct{}{
	fiber.MethodPost:   {},
	fiber.MethodPut:    {},
	fiber.MethodPatch:  {},
	fiber.MethodDelete: {},
}

func unsafeMethod(method string) bool {
	_, ok := unsafeMethods[method]
	return ok
}
