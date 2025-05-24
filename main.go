package main

import (
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	handler2 "notes-vault-backend/internal/handler"
	"notes-vault-backend/internal/middleware"
)

func main() {
	app := fiber.New()
	db, err := gorm.Open(sqlite.Open("./internal/db/notesapp.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	app.Post("/login", func(c *fiber.Ctx) error {
		return handler2.LoginHandler(c, db)
	})
	app.Post("/signup", middleware.JsonContentTypeValidator, middleware.SignupValidator, func(ctx *fiber.Ctx) error {
		return handler2.SignupHandler(ctx, db)
	})
	log.Fatal(app.Listen(":8080"))
}
