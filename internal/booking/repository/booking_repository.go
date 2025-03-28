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
	Create(ctx context.Context, booking *domain.Booking) error
	GetByID(ctx context.Context, id uint) (*domain.Booking, error)
	ListByProfessional(ctx context.Context, professionalID uint, from, to time.Time) ([]*domain.Booking, error)
	ListByClient(ctx context.Context, clientID uint) ([]*domain.Booking, error)
	IsTimeSlotAvailable(ctx context.Context, professionalID uint, start, end time.Time) (bool, error)
}

type bookingRepository struct {
	db *gorm.DB
}

func NewBookingRepository(db *gorm.DB) BookingRepository {
	return &bookingRepository{db: db}
}

func (r *bookingRepository) Create(ctx context.Context, booking *domain.Booking) error {
	// Verifica se o profissional existe
	var professionalExists bool
	if err := r.db.WithContext(ctx).
		Model(&professional.Professional{}).
		Select("count(*) > 0").
		Where("id = ?", booking.ProfessinalID).
		Find(&professionalExists).Error; err != nil {
		return err
	}
	if !professionalExists {
		return domain.ErrProfessionalUnavailable
	}

	// Verifica conflitos de horário
	if available, err := r.IsTimeSlotAvailable(ctx, booking.ProfessinalID, booking.StartTime, booking.EndTime); err != nil {
		return err
	} else if !available {
		return domain.ErrTimeSlotUnavailable
	}

	// Define valores padrão
	if booking.Status == "" {
		booking.Status = domain.StatusPending
	}
	booking.CreatedAt = time.Now()
	booking.UpdatedAt = time.Now()

	return r.db.WithContext(ctx).Create(booking).Error
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

func (r *bookingRepository) IsTimeSlotAvailable(ctx context.Context, professionalID uint, start, end time.Time) (bool, error) {
	var count int64

	err := r.db.WithContext(ctx).Model(&domain.Booking{}).Where("professional_id = ?", professionalID).Where("status IN ?", []domain.BookingStatus{
		domain.StatusPending,
		domain.StatusConfirmed,
	}).
		Where("(start_time, end_time) OVERLAPS (?, ?)", start, end).Count(&count).Error
	return count == 0, err
}
