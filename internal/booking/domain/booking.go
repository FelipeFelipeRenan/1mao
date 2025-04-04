package domain

import (
	"errors"
	"time"
)

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	StatusConfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
	StatusCompleted BookingStatus = "completed"
)

var (
	ErrBookingNotFound         = errors.New("booking not found")
	ErrTimeSlotUnavailable     = errors.New("time slot unavailable")
	ErrInvalidStatusTransition = errors.New("invalid status transition")
	ErrProfessionalUnavailable = errors.New("professional unavailable")
)

// Booking representa um usuário cliente do sistema
//
// Client representa um profissional
//
//	@Description	Modelo que representa o agendamento do profissinal para o cliente
//	@name			Booking
//	@model			Booking
type Booking struct {
	ID             uint          `json:"id" gorm:"primaryKey"`
	ProfessionalID uint          `json:"professional_id" gorm:"column:professional_id"` // Corrigido aqui
	ClientID       uint          `json:"client_id"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	Status         BookingStatus `json:"status"`
	CreatedAt      time.Time     `json:"created_at"`
	UpdatedAt      time.Time     `json:"updated_at"`
}

type Availability struct {
	ProfessionalID uint         `json:"professional_id" gorm:"primaryKey"`
	Weekday        time.Weekday `json:"weekday" gorm:"primaryKey"`
	StartHour      int          `json:"start_hour"` // 8 para as 08:00
	EndHour        int          `json:"end_hour"`   // 17 para 17:00
}
