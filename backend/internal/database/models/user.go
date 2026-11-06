package models

type User struct {
	Base
	Name         string `gorm:"not null" json:"name"`
	Email        string `gorm:"uniqueIndex;not null" json:"email"`
	PasswordHash string `gorm:"not null" json:"-"`
	Role         Role   `gorm:"not null;default:requester" json:"role"`
}
