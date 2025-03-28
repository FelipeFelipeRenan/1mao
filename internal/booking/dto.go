package booking

import (
	"1mao/internal/booking/domain"
	"time"
)

// DTOs temporários até criar o módulo service
type (
	CreateBookingRequest struct {
		ProfessionalID uint      `json:"professional_id"`
		ClientID      uint      `json:"client_id"`
		ServiceID     uint      `json:"service_id"`
		StartTime     time.Time `json:"start_time"`
		EndTime       time.Time `json:"end_time"`
	}

	BookingResponse struct {
		ID            uint          `json:"id"`
		ProfessionalID uint          `json:"professional_id"`
		ClientID      uint          `json:"client_id"`
		ServiceID     uint          `json:"service_id"`
		StartTime     time.Time     `json:"start_time"`
		EndTime       time.Time     `json:"end_time"`
		Status        domain.BookingStatus `json:"status"`
	}
)