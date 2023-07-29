package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/app"
	"github.com/tsivinsky/fileasy/internal/jwt"
)

func VerifyJWTToken(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	if authHeader == "" {
		return app.NewApiError(401, "No authorization header with token", nil)
	}

	values := strings.Split(authHeader, " ")
	tokenValue := values[1]

	_, err := jwt.ValidateAccessToken(string(tokenValue))
	if err != nil {
		return app.NewApiError(400, "Invalid accessToken", &err)
	}

	return c.Next()
}
