package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"notes-vault-backend/dto"
	"notes-vault-backend/model"
	"notes-vault-backend/utils"
	"strings"
)

func SignupHandler(ctx *fiber.Ctx, db *gorm.DB, err error) error {
	if err != nil {
		log.Fatal(err)
		return ctx.SendStatus(fiber.StatusInternalServerError)
	}
	contentType := ctx.Get("Content-Type")
	if contentType != "application/json" {
		return fiber.ErrUnsupportedMediaType
	}
	u := new(dto.SignUpRequest)
	if err := ctx.BodyParser(u); err != nil {
		return fiber.ErrBadRequest
	}
	u.Email = strings.Trim(u.Email, " ")
	u.Password = strings.Trim(u.Password, " ")
	if u.Email == "" || u.Password == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Email or Password is empty")
	}
	if utils.EmailValidator(u.Email) == false {
		return fiber.NewError(fiber.StatusBadRequest, "Invalid Email")
	}
	if len(u.Password) < 8 {
		return fiber.NewError(fiber.StatusBadRequest, "Password must be at least 8 characters")
	}
	result := db.First(&u, "email = ?", u.Email)
	if result.RowsAffected != 0 {
		return fiber.NewError(fiber.StatusConflict, "Email already exists")
	}
	entryUuid := uuid.Must(uuid.NewRandom()).String()
	entryHashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	databaseEntry := model.User{Id: entryUuid, Email: u.Email, HashedPassword: string(entryHashedPassword)}
	result = db.Select("id", "email", "hashed_password").Create(&databaseEntry)
	if result.Error != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	token, err := utils.JwtGenerator(databaseEntry)
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	return ctx.Status(201).JSON(fiber.Map{
		"token": token,
	})
}
