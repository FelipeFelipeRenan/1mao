package domain

import "gorm.io/gorm"


type AdminUser struct {
	gorm.Model
	Name string `gorm:"not null"`
	Email string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
}