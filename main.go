package main

import (
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"log"
	"notes-vault-backend/handler"
)

func main() {
	app := fiber.New()
	db, err := gorm.Open(sqlite.Open("./db/notesapp.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	app.Post("/login", func(c *fiber.Ctx) error {
		return handler.LoginHandler(c, db)
	})
	app.Post("/signup", func(ctx *fiber.Ctx) error {
		return handler.SignupHandler(ctx, db)
	})
	log.Fatal(app.Listen(":8080"))
}
