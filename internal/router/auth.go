package router

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/app"

	"github.com/tsivinsky/fileasy/internal/jwt"
)

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

type GetNewAccessTokenBody struct {
	RefreshToken string `json:"refreshToken"`
}

func HandleGetNewAccessToken(c *fiber.Ctx) error {
	var body GetNewAccessTokenBody

	if err := c.BodyParser(&body); err != nil {
		return app.NewApiError(400, "couldn't parse request body", &err)
	}

	if body.RefreshToken == "" {
		return app.NewApiError(400, "No refreshToken provided in body", nil)
	}

	userId, err := jwt.ValidateRefreshToken(body.RefreshToken)
	if err != nil {
		return app.NewApiError(400, "Invalid refreshToken", &err)
	}

	accessToken, refreshToken, err := jwt.GenerateBothTokens(userId)
	if err != nil {
		return app.NewApiError(500, "couldn't generate new tokens", &err)
	}

	return c.JSON(AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
