package main

import (
	"blenioviva/internal/db"
	tagshandler "blenioviva/internal/handlers/tags"

	"github.com/gofiber/fiber/v3"
)

func main() {
	app := fiber.New()

	db := db.New()

	tagshandler.AssignRoutes(app, db)

	app.Listen(":3000")
}
