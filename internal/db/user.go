package db

type User struct {
	Model

	Email    *string `json:"email" gorm:"email"`
	Username string  `json:"username" gorm:"username,unique"`
	Files    []File  `json:"files" gorm:"files"`
	YandexId *int    `json:"-" gorm:"yandex_id"`
	GithubId *int    `json:"-"`
}
