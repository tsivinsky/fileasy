package github

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GitHubUser struct {
	Email *string `json:"email"`
	Login string  `json:"login"`
}

func GetUserData(accessToken string) (*GitHubUser, error) {
	uReq, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return nil, err
	}
	uReq.Header.Set("Authorization", fmt.Sprintf("Bearer %s", accessToken))
	uReq.Header.Set("Accept", "application/json")

	uResp, err := http.DefaultClient.Do(uReq)
	if err != nil {
		return nil, err
	}
	defer uResp.Body.Close()

	ur, err := io.ReadAll(uResp.Body)
	if err != nil {
		return nil, err
	}

	var ghUser *GitHubUser
	err = json.Unmarshal(ur, &ghUser)
	if err != nil {
		return nil, err
	}

	return ghUser, nil
}
