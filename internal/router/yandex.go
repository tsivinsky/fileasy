package router

import (
	"fmt"
	"os"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/tsivinsky/fileasy/internal/db"
	"github.com/tsivinsky/fileasy/internal/jwt"
	"github.com/tsivinsky/fileasy/internal/yandex"
)

func HandleYandexLogin(c *fiber.Ctx) error {
	yandexClientId := os.Getenv("YANDEX_CLIENT_ID")
	fmt.Printf("yandexClientId: %v\n", yandexClientId)

	u := fmt.Sprintf("https://oauth.yandex.ru/authorize?response_type=code&client_id=%s", yandexClientId)

	return c.Redirect(u)
}

func HandleYandexCallback(c *fiber.Ctx) error {
	code := c.Query("code")

	accessToken, err := yandex.GetOauthToken(code)
	if err != nil {
		return err
	}

	yandexUser, err := yandex.GetYandexUser(accessToken)
	if err != nil {
		return err
	}

	yandexId, err := strconv.Atoi(yandexUser.ID)
	if err != nil {
		return err
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

		user.YandexId = &yandexId
		db.Db.Save(&user)
	}

	accessToken, refreshToken, err := jwt.GenerateBothTokens(user.ID)
	if err != nil {
		return err
	}

	return RedirectWithTokens(c, accessToken, refreshToken)
}
