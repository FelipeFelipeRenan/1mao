package domain

import "gorm.io/gorm"

type Professional struct {
	gorm.Model
	Name       string  `json:"name" gorm:"not null"`
	Email      string  `json:"email" gorm:"unique;not null"`
	Password   string  `json:"-" gorm:"not null"`
	Role       string  `json:"role" gorm:"default:professional"`
	Profession string  `json:"profession" gorm:"not null"`
	Experience int     `json:"experience"`
	Rating     float32 `json:"rating" gorm:"default:0"`
}
