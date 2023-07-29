package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/db"
)

func HandleGetUser(c *fiber.Ctx) error {
	userId, err := GetUserIdFromRequest(c)
	if err != nil {
		return err
	}

	var user db.User
	if tx := db.Db.Preload("Files").First(&user, "id = ?", userId); tx.Error != nil {
		return err
	}

	return c.JSON(user)
}
