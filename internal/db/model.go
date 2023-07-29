package db

import "time"

type Model struct {
	ID        uint       `json:"id" gorm:"primarykey"`
	CreatedAt time.Time  `json:"createdAt" gorm:"created_at"`
	UpdatedAt time.Time  `json:"updatedAt" gorm:"updated_at"`
	DeletedAt *time.Time `json:"deletedAt" gorm:"index"`
}
