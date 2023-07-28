package db

type User struct {
	Model

	Email    *string `json:"email" gorm:"email"`
	Username string  `json:"username" gorm:"username,unique"`
}