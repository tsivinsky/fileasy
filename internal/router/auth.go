package router

import (
	"errors"

	"github.com/gofiber/fiber/v2"
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
		return err
	}

	if body.RefreshToken == "" {
		return errors.New("No refreshToken provided in body")
	}

	userId, err := jwt.ValidateRefreshToken(body.RefreshToken)
	if err != nil {
		return err
	}

	accessToken, refreshToken, err := jwt.GenerateBothTokens(userId)
	if err != nil {
		return err
	}

	return c.JSON(AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}
