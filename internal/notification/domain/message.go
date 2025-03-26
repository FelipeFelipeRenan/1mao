// internal/notification/domain/message.go
package domain

import (
	"gorm.io/gorm"
	"time"
)

// Message representa uma mensagem no chat
// @Description Estrutura completa de mensagem com metadados
type Message struct {
	ID        uint           `gorm:"primaryKey" json:"id" example:"1"`
	CreatedAt time.Time      `json:"created_at" example:"2023-01-01T00:00:00Z"`
	UpdatedAt time.Time      `json:"updated_at" example:"2023-01-01T00:00:00Z"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-" swaggerignore:"true"`

	SenderID     int       `json:"sender_id" example:"1"`
	SenderType   string    `json:"sender_type" example:"client" enums:"client,professional"`
	ReceiverID   int       `json:"receiver_id" example:"2"`
	ReceiverType string    `json:"receiver_type" example:"professional" enums:"client,professional"`
	Content      string    `json:"content" example:"Olá, como posso ajudar?"`
	Timestamp    time.Time `json:"timestamp" example:"2023-01-01T15:04:05Z"`
}
