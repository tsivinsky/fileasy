package db

type User struct {
	Model

	Email    *string `json:"email"`
	Username string  `json:"username" gorm:"unique"`
	Files    []File  `json:"files"`
	GithubId *int    `json:"-"`
	YandexId *int    `json:"-"`
}
