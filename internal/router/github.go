package router

import (
	"fmt"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/app"
	"github.com/tsivinsky/fileasy/internal/db"
	"github.com/tsivinsky/fileasy/internal/github"
	"github.com/tsivinsky/fileasy/internal/jwt"
)

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
		return app.NewApiError(500, "couldn't get access token from github", &err)
	}

	ghUser, err := github.GetUserData(accessToken)
	if err != nil {
		return app.NewApiError(500, "couldn't get user data from github using access token", &err)
	}

	var user *db.User
	if tx := db.Db.Where("github_id = ? OR email = ?", ghUser.ID, ghUser.Email).First(&user); tx.Error != nil {
		user.Username = ghUser.Login
		user.Email = ghUser.Email
		user.GithubId = &ghUser.ID
		db.Db.Create(&user)
	} else {
		if user.Email == nil {
			user.Email = ghUser.Email
		}

		user.Avatar = ghUser.Avatar
		user.GithubId = &ghUser.ID
		db.Db.Save(&user)
	}

	accessToken, refreshToken, err := jwt.GenerateBothTokens(user.ID)
	if err != nil {
		return app.NewApiError(500, "couldn't generate jwt tokens", &err)
	}

	return RedirectWithTokens(c, accessToken, refreshToken)
}
