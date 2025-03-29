package service_test

import (
	"context"
	"testing"
	"time"

	"1mao/internal/booking/domain"
	"1mao/internal/booking/repository"
	"1mao/internal/booking/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockBookingRepository struct {
	mock.Mock
}

func (m *MockBookingRepository) Create(ctx context.Context, req *repository.CreateBookingRequest) (*domain.Booking, error) {
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

func (m *MockBookingRepository) ListByClient(ctx context.Context, clientID uint, from, to time.Time) ([]*domain.Booking, error) {
	args := m.Called(ctx, clientID, from, to)
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

func TestBookingService_ListClientBookings(t *testing.T) {
    ctx := context.Background()
    mockRepo := new(MockBookingRepository)
    bookingService := service.NewBookingService(mockRepo)

    // Mock data
    mockBookings := []*domain.Booking{
        {ID: 1, Status: domain.StatusConfirmed},
        {ID: 2, Status: domain.StatusCompleted},
    }

    t.Run("Success - all bookings", func(t *testing.T) {
        mockRepo.On("ListByClient", ctx, uint(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
            Return(mockBookings, nil).Once()

        // Chamada com filtros vazios
        result, err := bookingService.ListClientBookings(ctx, 1, &service.BookingFilters{})
        
        assert.NoError(t, err)
        assert.Equal(t, 2, len(result))
        mockRepo.AssertExpectations(t)
    })

    t.Run("Success - with date filters", func(t *testing.T) {
        from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
        to := time.Date(2025, 12, 31, 23, 59, 59, 0, time.UTC)
        
        mockRepo.On("ListByClient", ctx, uint(1), from, to).
            Return(mockBookings, nil).Once()

        // Chamada com filtros de data
        result, err := bookingService.ListClientBookings(ctx, 1, &service.BookingFilters{
            From: from,
            To:   to,
        })
        
        assert.NoError(t, err)
        assert.Equal(t, 2, len(result))
        mockRepo.AssertExpectations(t)
    })

    t.Run("Success - with status filter", func(t *testing.T) {
        mockRepo.On("ListByClient", ctx, uint(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
            Return(mockBookings, nil).Once()

        // Chamada com filtro de status
        result, err := bookingService.ListClientBookings(ctx, 1, &service.BookingFilters{
            Status: domain.StatusConfirmed,
        })
        
        assert.NoError(t, err)
        assert.Equal(t, 1, len(result))
        assert.Equal(t, domain.StatusConfirmed, result[0].Status)
        mockRepo.AssertExpectations(t)
    })

    t.Run("Error - repository error", func(t *testing.T) {
        mockRepo.On("ListByClient", ctx, uint(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
            Return(nil, assert.AnError).Once()

        // Chamada que deve retornar erro
        result, err := bookingService.ListClientBookings(ctx, 1, &service.BookingFilters{})
        
        assert.Error(t, err)
        assert.Nil(t, result)
        assert.Equal(t, assert.AnError, err)
        mockRepo.AssertExpectations(t)
    })
}
func TestBookingService_ListProfessionalBookings(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockBookingRepository)
	bookingService := service.NewBookingService(mockRepo)

	// Mock data
	mockBookings := []*domain.Booking{
		{ID: 1, Status: domain.StatusConfirmed},
		{ID: 2, Status: domain.StatusCompleted},
	}

	t.Run("Success - all bookings", func(t *testing.T) {
		mockRepo.On("ListByProfessional", ctx, uint(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return(mockBookings, nil).Once()

		result, err := bookingService.ListProfessionalBookings(ctx, 1, &service.BookingFilters{})
		
		assert.NoError(t, err)
		assert.Equal(t, 2, len(result))
		mockRepo.AssertExpectations(t)
	})

	t.Run("Success - with status filter", func(t *testing.T) {
		mockRepo.On("ListByProfessional", ctx, uint(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return(mockBookings, nil).Once()

		result, err := bookingService.ListProfessionalBookings(ctx, 1, &service.BookingFilters{
			Status: domain.StatusConfirmed,
		})
		
		assert.NoError(t, err)
		assert.Equal(t, 1, len(result))
		assert.Equal(t, domain.StatusConfirmed, result[0].Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - repository error", func(t *testing.T) {
		mockRepo.On("ListByProfessional", ctx, uint(1), mock.AnythingOfType("time.Time"), mock.AnythingOfType("time.Time")).
			Return(nil, assert.AnError).Once()

		result, err := bookingService.ListProfessionalBookings(ctx, 1, &service.BookingFilters{})
		
		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingService_CreateBooking(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockBookingRepository)
	bookingService := service.NewBookingService(mockRepo)

	futureTime := time.Now().Add(24 * time.Hour)

	t.Run("Success - create booking", func(t *testing.T) {
		req := &service.CreateBookingRequest{
			ProfessionalID: 1,
			ClientID:       2,
			StartTime:     futureTime,
			EndTime:       futureTime.Add(time.Hour),
		}

		expectedBooking := &domain.Booking{
			ID:             1,
			ProfessionalID: req.ProfessionalID,
			ClientID:       req.ClientID,
			StartTime:     req.StartTime,
			EndTime:       req.EndTime,
			Status:        domain.StatusPending,
		}

		mockRepo.On("IsTimeSlotAvailable", ctx, req.ProfessionalID, req.StartTime, req.EndTime).Return(true, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*repository.CreateBookingRequest")).Return(expectedBooking, nil).Once()

		result, err := bookingService.CreateBooking(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, expectedBooking.ID, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - time slot not available", func(t *testing.T) {
		req := &service.CreateBookingRequest{
			ProfessionalID: 1,
			ClientID:       2,
			StartTime:     futureTime,
			EndTime:       futureTime.Add(time.Hour),
		}

		mockRepo.On("IsTimeSlotAvailable", ctx, req.ProfessionalID, req.StartTime, req.EndTime).Return(false, nil).Once()

		result, err := bookingService.CreateBooking(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrTimeSlotUnavailable, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - repository error", func(t *testing.T) {
		req := &service.CreateBookingRequest{
			ProfessionalID: 1,
			ClientID:       2,
			StartTime:     futureTime,
			EndTime:       futureTime.Add(time.Hour),
		}

		mockRepo.On("IsTimeSlotAvailable", ctx, req.ProfessionalID, req.StartTime, req.EndTime).Return(true, nil).Once()
		mockRepo.On("Create", ctx, mock.AnythingOfType("*repository.CreateBookingRequest")).Return(nil, assert.AnError).Once()

		result, err := bookingService.CreateBooking(ctx, req)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingService_GetBooking(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockBookingRepository)
	bookingService := service.NewBookingService(mockRepo)

	t.Run("Success - get booking", func(t *testing.T) {
		bookingID := uint(1)
		expectedBooking := &domain.Booking{
			ID:     bookingID,
			Status: domain.StatusConfirmed,
		}

		mockRepo.On("GetByID", ctx, bookingID).Return(expectedBooking, nil).Once()

		result, err := bookingService.GetBooking(ctx, bookingID)

		assert.NoError(t, err)
		assert.Equal(t, expectedBooking.ID, result.ID)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - booking not found", func(t *testing.T) {
		bookingID := uint(999)

		mockRepo.On("GetByID", ctx, bookingID).Return(nil, domain.ErrBookingNotFound).Once()

		result, err := bookingService.GetBooking(ctx, bookingID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrBookingNotFound, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - repository error", func(t *testing.T) {
		bookingID := uint(1)

		mockRepo.On("GetByID", ctx, bookingID).Return(nil, assert.AnError).Once()

		result, err := bookingService.GetBooking(ctx, bookingID)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}

func TestBookingService_UpdateBookingStatus(t *testing.T) {
	ctx := context.Background()
	mockRepo := new(MockBookingRepository)
	bookingService := service.NewBookingService(mockRepo)

	t.Run("Success - update status", func(t *testing.T) {
		bookingID := uint(1)
		newStatus := domain.StatusConfirmed
		expectedBooking := &domain.Booking{
			ID:     bookingID,
			Status: newStatus,
		}

		mockRepo.On("UpdateStatus", ctx, bookingID, newStatus).Return(expectedBooking, nil).Once()

		result, err := bookingService.UpdateBookingStatus(ctx, bookingID, newStatus)

		assert.NoError(t, err)
		assert.Equal(t, expectedBooking.ID, result.ID)
		assert.Equal(t, expectedBooking.Status, result.Status)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - invalid status transition", func(t *testing.T) {
		bookingID := uint(1)
		newStatus := domain.StatusCompleted

		mockRepo.On("UpdateStatus", ctx, bookingID, newStatus).Return(nil, domain.ErrInvalidStatusTransition).Once()

		result, err := bookingService.UpdateBookingStatus(ctx, bookingID, newStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, domain.ErrInvalidStatusTransition, err)
		mockRepo.AssertExpectations(t)
	})

	t.Run("Error - repository error", func(t *testing.T) {
		bookingID := uint(1)
		newStatus := domain.StatusConfirmed

		mockRepo.On("UpdateStatus", ctx, bookingID, newStatus).Return(nil, assert.AnError).Once()

		result, err := bookingService.UpdateBookingStatus(ctx, bookingID, newStatus)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, assert.AnError, err)
		mockRepo.AssertExpectations(t)
	})
}