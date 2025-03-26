package repository

import (
	"1mao/internal/notification/domain"
	"fmt"
	"gorm.io/gorm"
)

// Estrutura do repositório de mensagens
type MessageRepository struct {
	DB *gorm.DB
}

// Criar um novo repositório de mensagens
func NewMessageRepository(db *gorm.DB) *MessageRepository {
	return &MessageRepository{DB: db}
}

// Salvar mensagem no banco de dados
func (repo *MessageRepository) SaveMessage(msg domain.Message) error {
	msg.Timestamp = msg.Timestamp.UTC()
	result := repo.DB.Create(&msg)
	if result.Error != nil {
		return result.Error
	}
	fmt.Printf("💾 Mensagem salva -> %s\n", msg.Content)
	return nil
}

// Recuperar mensagens antigas entre dois usuários
func (repo *MessageRepository) GetMessages(userID int, userType string, otherID int, otherType string) ([]domain.Message, error) {
	var messages []domain.Message

	err := repo.DB.Where(
		"(sender_id = ? AND sender_type = ? AND receiver_id = ? AND receiver_type = ?) OR "+
			"(sender_id = ? AND sender_type = ? AND receiver_id = ? AND receiver_type = ?)",
		userID, userType, otherID, otherType,
		otherID, otherType, userID, userType,
	).Order("timestamp ASC").Find(&messages).Error

	if err != nil {
		return nil, err
	}

	return messages, nil
}
// Testes de conexão de websocket
// {"sender_id": 1, "sender_type": "professional", "receiver_id": 3,"receiver_type": "client","content": "Olá!" }
// {"sender_id": 3, "sender_type": "client", "receiver_id": 1,"receiver_type": "professional","content": "Tudo bem?" }