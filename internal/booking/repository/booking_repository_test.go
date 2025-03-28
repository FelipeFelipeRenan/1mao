package repository

import (
	"1mao/internal/booking/domain"
	"context"
	"time"

	"github.com/stretchr/testify/mock"
)

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) CreateBooking(ctx context.Context, req *CreateBookingRequest) (*domain.Booking, error) {
	args := m.Called(ctx, req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) GetByID(ctx context.Context, id uint) (*domain.Booking, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) ListByProfessional(ctx context.Context, professionalID uint, from, to time.Time) ([]*domain.Booking, error) {
	args := m.Called(ctx, professionalID, from, to)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) ListByClient(ctx context.Context, clientID uint) ([]*domain.Booking, error) {
	args := m.Called(ctx, clientID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) UpdateStatus(ctx context.Context, id uint, status domain.BookingStatus) (*domain.Booking, error) {
	args := m.Called(ctx, id, status)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Booking), args.Error(1)
}

func (m *MockBookingRepository) IsTimeSlotAvailable(ctx context.Context, professionalID uint, start, end time.Time) (bool, error) {
	args := m.Called(ctx, professionalID, start, end)
	return args.Bool(0), args.Error(1)
}