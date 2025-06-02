package tagshandler

import (
	"blenioviva/internal/db"

	"github.com/gofiber/fiber/v3"
)

type TagsHandler struct {
	db *db.DB
}

func newTagsHandler(db *db.DB) *TagsHandler {
	return &TagsHandler{
		db: db,
	}
}

func AssignRoutes(app *fiber.App, db *db.DB) {
	handlers := newTagsHandler(db)

	app.Post("/tags", handlers.CreateTag)
}
