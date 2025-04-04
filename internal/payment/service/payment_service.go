package service

import (
	"1mao/internal/payment/domain"
	"1mao/internal/payment/repository"

	"github.com/google/uuid"
)

type PaymentService interface {
	CreatePayment(clientID string, bookingID string, amount int64, method string) (*domain.Transaction, error)
	ConfirmPayment(gatewayID string) error
	FailPayment(gatewayID string) error
	GetPaymentByID(paymentID string) (*domain.Transaction, error)
	GetClientPayments(clientID string) ([]domain.Transaction, error)
}

type paymentService struct {
	repo   repository.PaymentRepository
	stripe *StripeClient
}

func NewPaymentService(repo repository.PaymentRepository, stripeKey string) PaymentService {
	return &paymentService{
		repo:   repo,
		stripe: NewStripeClient(stripeKey),
	}
}

func (s *paymentService) CreatePayment(clientID string, bookingID string, amount int64, method string) (*domain.Transaction, error) {
	// criar intent no stripe
	intent, err := s.stripe.CreatePaymentIntent(amount, "brl")
	if err != nil {
		return nil, err
	}

	transaction := domain.Transaction{
		ID:            uuid.NewString(),
		BookingID:     bookingID,
		ClientID:      clientID,
		Amount:        amount,
		Currency:      "BRL",
		Status:        domain.StatusPending,
		PaymentMethod: method,
		GatewayID:     intent.ID,
	}

	err = s.repo.CreateTransaction(transaction)
	return &transaction, err
}

func (s *paymentService) ConfirmPayment(gatewayID string) error {
	return s.repo.UpdateStatus(gatewayID, string(domain.StatusPaid))
}

func (s *paymentService) FailPayment(gatewayID string) error {
	return s.repo.UpdateStatus(gatewayID, string(domain.StatusFailed))
}

func (s *paymentService) GetPaymentByID(paymentID string) (*domain.Transaction, error) {
	return s.repo.GetByID(paymentID)
}

func (s *paymentService) GetClientPayments(clientID string) ([]domain.Transaction, error) {
	return s.repo.GetByClientID(clientID)
}
