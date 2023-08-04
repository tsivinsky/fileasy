package db

type File struct {
	Model

	Name   string `json:"name" gorm:"unique"`
	User   User   `json:"-"`
	UserID uint   `json:"userId"`
}
