package domain

import (
	"gorm.io/gorm"
	"time"
)

// Estrutura da mensagem do chat no banco de dados
type Message struct {
	gorm.Model
	ID           uint      `gorm:"primaryKey"`
	SenderID     int       `json:"sender_id"`
	SenderType   string    `json:"sender_type"`
	ReceiverID   int       `json:"receiver_id"`
	ReceiverType string    `json:"receiver_type"`
	Content      string    `json:"content"`
	Timestamp    time.Time `json:"timestamp"`
}
