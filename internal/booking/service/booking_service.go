package service

import (
	"1mao/internal/booking/repository"
	"time"
)

type BookingService interface {
}

type bookingService struct {
	bookingRepo *repository.BookingRepository
}

type CreateBookingRequest struct {
	ProfessionalID uint `json:"professional_id"`
	ClientID uint `json:"client_id"`
	ServiceID uint `json:"service_id"`
	Date time.Time `json:"date"`
	StartTime time.Time `json:"start_time"`
	EndTime time.Time `json:"end_time"`
}

func NewBookingService(bookingRepo repository.BookingRepository) BookingService {
	return &bookingService{bookingRepo: &bookingRepo}
}

// TODO: Metodos do service