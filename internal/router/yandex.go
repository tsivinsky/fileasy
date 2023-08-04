package router

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/app"
	"github.com/tsivinsky/fileasy/internal/db"
	"github.com/tsivinsky/fileasy/internal/jwt"
	"github.com/tsivinsky/fileasy/internal/yandex"
)

const YandexAvatarSize = "islands-middle"

func HandleYandexLogin(c *fiber.Ctx) error {
	yandexClientId := os.Getenv("YANDEX_CLIENT_ID")

	u := fmt.Sprintf("https://oauth.yandex.ru/authorize?response_type=code&client_id=%s", yandexClientId)

	return c.Redirect(u)
}

func HandleYandexCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	accessToken, err := yandex.GetOauthToken(code)
	if err != nil {
		return app.NewApiError(500, "couldn't get yandex oauth token", &err)
	}

	yandexUser, err := yandex.GetYandexUser(accessToken)
	if err != nil {
		return app.NewApiError(500, "couldn't get yandex user", &err)
	}

	yandexId, err := strconv.Atoi(yandexUser.ID)
	if err != nil {
		return app.NewApiError(500, "couldn't convert yandex userId from string to int", &err)
	}

	var user db.User
	if tx := db.Db.First(&user, "yandex_id = ? OR email = ?", yandexId, yandexUser.Emails[0]); tx.Error != nil {
		user.Username = yandexUser.Login
		user.Email = &yandexUser.Emails[0]
		user.YandexId = &yandexId
		db.Db.Create(&user)
	} else {
		if user.Email == nil {
			user.Email = &yandexUser.Emails[0]
		}

		avatarUrl := fmt.Sprintf("https://avatars.yandex.net/get-yapic/%s/%s", yandexUser.DefaultAvatarId, YandexAvatarSize)
		user.Avatar = &avatarUrl

		user.YandexId = &yandexId
		db.Db.Save(&user)
	}

	accessToken, refreshToken, err := jwt.GenerateBothTokens(user.ID)
	if err != nil {
		return app.NewApiError(500, "couldn't generate jwt tokens", &err)
	}

	return RedirectWithTokens(c, accessToken, refreshToken)
}
