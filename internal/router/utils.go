package router

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/jwt"
)

func GetUserIdFromRequest(c *fiber.Ctx) (uint, error) {
	authHeader := c.Get("Authorization")
	values := strings.Split(authHeader, " ")
	accessToken := values[1]

	userId, err := jwt.ValidateAccessToken(accessToken)
	if err != nil {
		return 0, err
	}

	return userId, nil
}
