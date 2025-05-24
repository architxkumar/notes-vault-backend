package middleware

import (
	"github.com/gofiber/fiber/v2"
	"notes-vault-backend/internal/dto"
	"notes-vault-backend/internal/utils"
	"strings"
)

func SignupValidator(c *fiber.Ctx) error {
	body := new(dto.SignUpRequest)
	if err := c.BodyParser(body); err != nil {
		return fiber.ErrBadRequest
	}
	email := strings.TrimSpace(body.Email)
	password := strings.TrimSpace(body.Password)
	if email == "" || password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email or Password is empty")
	}
	if utils.EmailValidator(email) == false {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Email")
	}
	if len(password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "Password length should be at least 8 characters")
	}
	c.Locals("signup_payload", dto.SignUpRequest{
		Email:    email,
		Password: password,
	})
	return c.Next()
}
