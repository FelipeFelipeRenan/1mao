// internal/client/domain/client.go
package domain

import (
	"gorm.io/gorm"
	"time"
)

type Role string

const (
	RoleClient       Role = "client"
	RoleProfessional Role = "professional"
)

// Client representa um usuário cliente do sistema
//	@name	Client
//	@model	Client
type Client struct {
	gorm.Model

	//	@Example	1
	ID uint `json:"id" gorm:"primaryKey"`

	//	@Example	João Cliente
	Name string `json:"name" gorm:"not null"`

	//	@Example	cliente@email.com
	Email string `json:"email" gorm:"unique;not null"`

	Password string `json:"-" gorm:"not null" swaggerignore:"true"`

	//	@Enum	client,professional
	Role Role `json:"role" gorm:"type:varchar(20);not null;default:client"`

	//	@Example	2023-01-01T00:00:00Z
	LastLogin time.Time `json:"last_login"`

	//	@Example	+5511999999999
	Phone string `json:"phone"`

	ResetToken string `json:"-" swaggerignore:"true"`

	ResetTokenExpiry time.Time `json:"-" swaggerignore:"true"`
}
