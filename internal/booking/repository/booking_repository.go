package repository

import (
	"1mao/internal/booking/domain"
	professional "1mao/internal/professional/domain"
	"errors"
	"time"

	"context"

	"gorm.io/gorm"
)

type BookingRepository interface {
	Create(ctx context.Context, req *CreateBookingRequest) (*domain.Booking, error)
	GetByID(ctx context.Context, id uint) (*domain.Booking, error)
	ListByProfessional(ctx context.Context, professionalID uint, from, to time.Time) ([]*domain.Booking, error)
	ListByClient(ctx context.Context, clientID uint) ([]*domain.Booking, error)
	UpdateStatus(ctx context.Context, id uint, status domain.BookingStatus) (*domain.Booking, error)
	IsTimeSlotAvailable(ctx context.Context, professionalID uint, start, end time.Time) (bool, error)
}

// CreateBookingRequest - Temporário até criar o módulo service
type CreateBookingRequest struct {
	ProfessionalID uint      `json:"professional_id"`
	ClientID       uint      `json:"client_id"`
	ServiceID      uint      `json:"service_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}


func (r *bookingRepository) Create(ctx context.Context, req *CreateBookingRequest) (*domain.Booking, error) {
	// Verifica se o profissional existe
	var professionalExists bool
	if err := r.db.WithContext(ctx).
		Model(&professional.Professional{}).
		Select("count(*) > 0").
		Where("id = ?", req.ProfessionalID).
		Find(&professionalExists).Error; err != nil {
		return nil, err
	}
	if !professionalExists {
		return nil, domain.ErrProfessionalUnavailable
	}

	// Verifica conflitos de horário
	if available, err := r.IsTimeSlotAvailable(ctx, req.ProfessionalID, req.StartTime, req.EndTime); err != nil {
		return nil, err
	} else if !available {
		return nil, domain.ErrTimeSlotUnavailable
	}

	// Cria o booking
	booking := &domain.Booking{
		ProfessinalID: req.ProfessionalID,
		ClientID:      req.ClientID,
		ServiceID:     req.ServiceID,
		StartTime:     req.StartTime,
		EndTime:       req.EndTime,
		Status:        domain.StatusPending,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	err := r.db.WithContext(ctx).Create(booking).Error
	return booking, err
}

func (r *bookingRepository) GetByID(ctx context.Context, id uint) (*domain.Booking, error) {
	var booking domain.Booking
	err := r.db.WithContext(ctx).
		Preload("Professional").
		Preload("Client").
		First(&booking, id).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrBookingNotFound
		}
		return nil, err
	}
	return &booking, nil
}

func (r *bookingRepository) ListByProfessional(ctx context.Context, professionalID uint, from, to time.Time) ([]*domain.Booking, error) {
	var bookings []*domain.Booking

	query := r.db.WithContext(ctx).
		Where("professional_id = ?", professionalID).
		Preload("Client").
		Preload("Service")

	if !from.IsZero() {
		query = query.Where("start_time >= ?", from)
	}
	if !to.IsZero() {
		query = query.Where("end_time <= ?", to)
	}

	err := query.Order("start_time ASC").Find(&bookings).Error
	return bookings, err

}

func (r *bookingRepository) ListByClient(ctx context.Context, clientID uint) ([]*domain.Booking, error) {
	var bookings []*domain.Booking
	err := r.db.WithContext(ctx).
		Where("client_id = ?", clientID).
		Preload("Professional").
		Preload("Service").
		Order("start_time DESC").
		Find(&bookings).Error
	return bookings, err
}

func (r *bookingRepository) UpdateStatus(ctx context.Context, id uint, status domain.BookingStatus) (*domain.Booking, error) {
	var booking domain.Booking

	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// Lock na transação
		if err := tx.Set("gorm:query_option", "FOR UPDATE").First(&booking, id).Error; err != nil {
			return err
		}

		// valida transição de status
		if !isValidStatusTransition(booking.Status, status) {
			return domain.ErrInvalidStatusTransition
		}

		booking.Status = status
		booking.UpdatedAt = time.Now()
		return tx.Save(&booking).Error

	})
	return &booking, err
}

func isValidStatusTransition(current, newStatus domain.BookingStatus) bool {
	transitions := map[domain.BookingStatus][]domain.BookingStatus{
		domain.StatusPending:   {domain.StatusConfirmed, domain.StatusCancelled},
		domain.StatusConfirmed: {domain.StatusCompleted, domain.StatusCancelled},
		domain.StatusCancelled: {},
		domain.StatusCompleted: {},
	}

	allowed, exists := transitions[current]
	if !exists {
		return false
	}

	for _, s := range allowed {
		if s == newStatus {
			return true
		}
	}
	return false
}

func (r *bookingRepository) IsTimeSlotAvailable(ctx context.Context, professionalID uint, start, end time.Time) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&domain.Booking{}).Where("professional_id = ?", professionalID).Where("status IN ?", []domain.BookingStatus{
		domain.StatusPending,
		domain.StatusConfirmed,
	}).
		Where("(start_time, end_time) OVERLAPS (?, ?)", start, end).Count(&count).Error
	return count == 0, err
}
