package handler

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"notes-vault-backend/dto"
	"notes-vault-backend/model"
	"notes-vault-backend/utils"
	"strings"
)

func LoginHandler(c *fiber.Ctx, db *gorm.DB, err error) error {
	if err != nil {
		log.Fatal(err)
		return fiber.ErrInternalServerError
	}
	if strings.Contains(c.Get("Content-type"), "application/json") == false {
		return fiber.ErrBadRequest
	}
	userCredential := new(dto.LoginRequest)
	err = c.BodyParser(userCredential)
	if err != nil {
		return fiber.ErrBadRequest
	}
	userCredential.Email = strings.TrimSpace(userCredential.Email)
	userCredential.Password = strings.TrimSpace(userCredential.Password)
	if userCredential.Email == "" || userCredential.Password == "" {
		return fiber.ErrBadRequest
	}
	if utils.EmailValidator(userCredential.Email) == false {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email address or password")
	}
	if len(userCredential.Password) < 8 {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email address or password")
	}
	var user model.User
	result := db.Where("email = ?", userCredential.Email).First(&user)
	if result.Error != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email address or password")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.HashedPassword), []byte(userCredential.Password))
	if err != nil {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid email address or password")
	}
	jwtToken, err := utils.JwtGenerator(user)
	if err != nil {
		return fiber.ErrInternalServerError
	}
	return c.JSON(fiber.Map{"token": jwtToken})
}
