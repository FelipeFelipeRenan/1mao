package domain

import "gorm.io/gorm"

// AdminUser representa um usuário administrador do sistema
// @Description Modelo completo de usuário administrador com controle total
type AdminUser struct {
	gorm.Model
	ID uint `json:"id" gorm:"primaryKey"`

	Name string `json:"name" gorm:"not null"`

	Email string `json:"email" gorm:"unique;not null"`

	Password string `json:"-" gorm:"not null"` // Oculto na resposta JSON

	IsActive bool `json:"is_active" gorm:"default:true"`
}
