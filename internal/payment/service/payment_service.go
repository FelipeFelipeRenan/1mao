package service

import (
	"1mao/internal/payment/domain"
	"1mao/internal/payment/repository"
)

type PaymentService interface {
	CreatePayment(bookingID string, amount int64, method string) (*domain.Transaction, error)
	ConfirmPayment(gatewayID string) error
	FailPayment(gatewayID string) error
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

func (s *paymentService) CreatePayment(bookingID string, amount int64, method string) (*domain.Transaction, error) {
	// criar intent no stripe
	intent, err := s.stripe.CreatePaymentIntent(amount, "brl")
	if err != nil {
		return nil, err
	}

	transaction := domain.Transaction{
		BookingID:     bookingID,
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
