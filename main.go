package main

import (
	"github.com/gofiber/fiber/v2"
	"notes-vault-backend/handlers"
)

func main() {
	app := fiber.New()
	app.Post("/signup", handlers.SignupHandler)
	app.Listen(":8080")
}
