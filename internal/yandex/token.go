package yandex

import (
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type OauthTokenResponse struct {
	AccessToken string `json:"access_token"`
}

func GetOauthToken(code string) (string, error) {
	clientId := os.Getenv("YANDEX_CLIENT_ID")
	clientSecret := os.Getenv("YANDEX_CLIENT_SECRET")

	form := url.Values{}
	form.Set("grant_type", "authorization_code")
	form.Set("code", code)
	form.Set("client_id", clientId)
	form.Set("client_secret", clientSecret)

	req, err := http.NewRequest("POST", "https://oauth.yandex.ru/token", strings.NewReader(form.Encode()))
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var result OauthTokenResponse
	if err := json.Unmarshal(data, &result); err != nil {
		return "", err
	}

	return result.AccessToken, nil
}
