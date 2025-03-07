package domain

import (
	"time"

	"gorm.io/gorm"
)

type Role string

const (
	RoleClient       Role = "client"
	RoleProfessional Role = "professional"
)

type User struct {
	gorm.Model
	Name     string `gorm:"not null"`
	Email    string `gorm:"unique;not null"`
	Password string `gorm:"not null"`
	Role     Role   `gorm:"type:varchar(20);not null"`

	ResetToken string 
	ResetTokenExpiry time.Time
}
