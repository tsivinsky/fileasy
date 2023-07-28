package github

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type AccessTokenBody struct {
	ClientId     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	Code         string `json:"code"`
}

type AccessTokenResult struct {
	AccessToken string `json:"access_token"`
}

func GetAccessToken(code, clientId, clientSecret string) (string, error) {
	ghTokenBody := AccessTokenBody{
		ClientId:     clientId,
		ClientSecret: clientSecret,
		Code:         code,
	}
	b, err := json.Marshal(&ghTokenBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", "https://github.com/login/oauth/access_token", bytes.NewReader(b))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	r, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var res AccessTokenResult
	err = json.Unmarshal(r, &res)
	if err != nil {
		return "", err
	}

	if res.AccessToken == "" {
		return "", errors.New("couldn't get access token from github")
	}

	return res.AccessToken, nil
}
