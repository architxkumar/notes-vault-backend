package middleware

import (
	"github.com/gofiber/fiber/v2"
	"strings"
)

func JsonContentTypeValidator(c *fiber.Ctx) error {
	if strings.Contains(strings.ToLower(c.Get("Content-Type")), "application/json") == false {
		return fiber.ErrUnsupportedMediaType
	}
	return c.Next()
}
