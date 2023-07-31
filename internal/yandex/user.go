package yandex

import (
	"encoding/json"
	"io"
	"net/http"
)

type YandexUser struct {
	ID              string   `json:"id"`
	Login           string   `json:"login"`
	Emails          []string `json:"emails"`
	IsAvatarEmpty   bool     `json:"is_avatar_empty"`
	DefaultAvatarId string   `json:"default_avatar_id"`
}

func GetYandexUser(accessToken string) (*YandexUser, error) {
	req, err := http.NewRequest("GET", "https://login.yandex.ru/info", nil)
	if err != nil {
		return nil, err
	}

	req.URL.Query().Set("format", "json")
	req.Header.Set("Authorization", "Oauth "+accessToken)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var yandexUser YandexUser
	if err := json.Unmarshal(data, &yandexUser); err != nil {
		return nil, err
	}

	return &yandexUser, nil
}
