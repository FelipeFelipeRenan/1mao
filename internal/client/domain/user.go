// internal/client/domain/client.go
package domain

import (
	_ "gorm.io/gorm"
	"time"
)

type Role string

const (
	RoleClient       Role = "client"
	RoleProfessional Role = "professional"
)

// Client representa um profissional
//
//	@Description	Modelo completo de cliente
//	@name			Client
//	@model			Client
type Client struct {
	ID               uint      `json:"id" gorm:"primaryKey"`
	Name             string    `json:"name" gorm:"not null"`
	Email            string    `json:"email" gorm:"unique;not null"`
	Password         string    `json:"-" gorm:"not null" swaggerignore:"true"`
	Role             Role      `json:"role" gorm:"type:varchar(20);not null;default:client"`
	LastLogin        time.Time `json:"last_login"`
	Phone            string    `json:"phone"`
	ResetToken       string    `json:"-" swaggerignore:"true"`
	ResetTokenExpiry time.Time `json:"-" swaggerignore:"true"`
}
