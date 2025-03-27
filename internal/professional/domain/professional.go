// internal/professional/domain/professional.go
package domain

import "time"

// Professional representa um profissional
// @Description Modelo completo de profissional
type Professional struct {
    ID        uint      `json:"id" gorm:"primaryKey"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
    
    Name       string  `json:"name" gorm:"not null"`
    Email      string  `json:"email" gorm:"unique;not null"`
    Password   string  `json:"-" gorm:"not null"`
    Profession string  `json:"profession" gorm:"not null"`
    Experience int     `json:"experience" gorm:"default:0"`
    Rating     float32 `json:"rating" gorm:"default:0"`
    Verified   bool    `json:"verified" gorm:"default:false"`
}