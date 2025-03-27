package domain

import "time"

type BookingStatus string

const (
	StatusPending   BookingStatus = "pending"
	Statusconfirmed BookingStatus = "confirmed"
	StatusCancelled BookingStatus = "cancelled"
	StatusCompleted BookingStatus = "completed"
)


// Booking representa um usuário cliente do sistema
//
// Client representa um profissional
//	@Description	Modelo que representa o agendamento do profissinal para o cliente
//	@name			Booking
//	@model			Booking
type Booking struct {
	ID            string        `json:"id" gorm:"primaryKey"`
	ProfessinalID string        `json:"professional_id"`
	ClientID      string        `json:"client_id"`
	ServiceID     string        `json:"service_id"`
	Date          time.Time     `json:"date"`
	StartTime     time.Time     `json:"start_time"`
	EndTime       time.Time     `json:"end_time"`
	Status        BookingStatus `json:"status"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

type Availability struct {
	ProfessionalID string       `json:"professional_id" gorm:"primaryKey"`
	Weekday        time.Weekday `json:"weekday" gorm:"primaryKey"`
	StartHour      int          `json:"start_hour"` // 8 para as 08:00
	EndHour        int          `json:"end_hour"`   // 17 para 17:00
}
