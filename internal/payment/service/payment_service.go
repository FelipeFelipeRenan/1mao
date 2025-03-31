package service

import "1mao/internal/payment/repository"

type PaymentService interface {
}
type paymentService struct {
	repo repository.PaymentRepository
}

func NewPaymentService(repo repository.PaymentRepository) PaymentService {
	return &paymentService{
		repo: repo,
	}
}
