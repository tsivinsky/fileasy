package middleware

import (
	"errors"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/jwt"
)

func VerifyJWTToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return errors.New("No authorization header with token")
	}

	values := strings.Split(authHeader, " ")
	tokenValue := values[1]

	_, err := jwt.ValidateAccessToken(string(tokenValue))
	if err != nil {
		return err
	}

	return c.Next()
}
