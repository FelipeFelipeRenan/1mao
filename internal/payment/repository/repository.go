package repository

import (
	"1mao/internal/payment/domain"

	"gorm.io/gorm"
)

type PaymentRepository interface {
	CreateTransaction(transaction domain.Transaction) error
	GetByGatewayID(gatewayID string) (*domain.Transaction, error)
	UpdateStatus(transactionID string, status string) error
}

type paymentRepository struct {
	db *gorm.DB
}

func NewPaymentRepository(db *gorm.DB) PaymentRepository {
	return &paymentRepository{db: db}
}

// CreateTransaction implements PaymentRepository.
func (p *paymentRepository) CreateTransaction(transaction domain.Transaction) error {
	return p.db.Create(transaction).Error
}

func (p *paymentRepository) GetByGatewayID(gatewayID string) (*domain.Transaction, error) {
	var transaction domain.Transaction
	err := p.db.Where("gateway_id = ?", gatewayID).First(&transaction).Error
	return &transaction, err
}

func (p *paymentRepository) UpdateStatus(transactionID string, status string) error {
	return p.db.Model(&domain.Transaction{}).
		Where("id = ?", transactionID).
		Update("status", status).Error
}
