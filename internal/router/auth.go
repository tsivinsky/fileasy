package router

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/db"
	"github.com/tsivinsky/fileasy/internal/github"
	"github.com/tsivinsky/fileasy/internal/jwt"
)

type AuthResponse struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
}

const SessionLifeTime = time.Second * 60 * 60 * 24 * 30

func HandleGitHubLogin(c *fiber.Ctx) error {
	clientId := os.Getenv("GITHUB_CLIENT_ID")

	uri := fmt.Sprintf("https://github.com/login/oauth/authorize?scope=read:user&client_id=%s", clientId)

	return c.Redirect(uri)
}

func HandleGitHubCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	clientId := os.Getenv("GITHUB_CLIENT_ID")
	clientSecret := os.Getenv("GITHUB_CLIENT_SECRET")

	accessToken, err := github.GetAccessToken(code, clientId, clientSecret)
	if err != nil {
		return err
	}

	ghUser, err := github.GetUserData(accessToken)
	if err != nil {
		return err
	}

	var user *db.User
	if tx := db.Db.Where("username = ?", ghUser.Login).First(&user); tx.Error != nil {
		user.Username = ghUser.Login
		user.Email = ghUser.Email
		db.Db.Create(&user)
	}

	accessToken, refreshToken, err := jwt.GenerateBothTokens(user.ID)
	if err != nil {
		return err
	}

	return c.Status(200).JSON(AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	})
}