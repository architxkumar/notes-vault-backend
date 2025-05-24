package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/log"
	"github.com/google/uuid"
	_ "github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"notes-vault-backend/internal/dto"
	"notes-vault-backend/internal/model"
	utils2 "notes-vault-backend/internal/utils"
)

func SignupHandler(ctx *fiber.Ctx, db *gorm.DB) error {
	u := ctx.Locals("signup_payload").(dto.SignUpRequest)
	var user model.User
	result := db.First(&user, "email = ?", u.Email)
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
	token, err := utils2.JwtGenerator(databaseEntry)
	if err != nil {
		log.Error(err)
		return fiber.ErrInternalServerError
	}
	return ctx.Status(201).JSON(fiber.Map{
		"token": token,
	})
}
