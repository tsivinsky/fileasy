package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/app"
	"github.com/tsivinsky/fileasy/internal/db"
)

func HandleGetUser(c *fiber.Ctx) error {
	userId, err := GetUserIdFromRequest(c)
	if err != nil {
		return app.NewApiError(401, "Unauthorized", &err)
	}

	var user db.User
	if tx := db.Db.Preload("Files").First(&user, "id = ?", userId); tx.Error != nil {
		return app.NewApiError(404, "couldn't find user by id", &tx.Error)
	}

	return c.JSON(user)
}
