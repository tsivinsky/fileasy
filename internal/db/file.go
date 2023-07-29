package db

type File struct {
	Model

	Name   string `json:"name" gorm:"name,unique"`
	UserID uint   `json:"userId" gorm:"user_id"`
}
