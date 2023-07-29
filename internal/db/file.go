package db

type File struct {
	Model

	Name   string `json:"name" gorm:"name,unique"`
	User   User   `json:"-" gorm:"user"`
	UserID uint   `json:"userId" gorm:"user_id"`
}
