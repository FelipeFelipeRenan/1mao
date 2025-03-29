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

func TestCreateBooking_Success(t *testing.T) {
	// Setup
	mockRepo := new(repository.MockBookingRepository)
	svc := service.NewBookingService(mockRepo)

	now := time.Now()
	startTime := now.Add(24 * time.Hour)
	endTime := startTime.Add(2 * time.Hour)

	req := &service.CreateBookingRequest{
		ProfessionalID: 1,
		ClientID:       1,
		StartTime:      startTime,
		EndTime:        endTime,
	}

	// Configuração do mock com tipos explícitos
	mockRepo.On("IsTimeSlotAvailable",
		mock.MatchedBy(func(ctx context.Context) bool { return true }), // Context
		uint(1),                          // professionalID (uint)
		mock.AnythingOfType("time.Time"), // start
		mock.AnythingOfType("time.Time"), // end
	).Return(true, nil)

	expectedBooking := &domain.Booking{
		ID:            1,
		ProfessionalID: 1,
		ClientID:      1,
		StartTime:     startTime,
		EndTime:       endTime,
		Status:        domain.StatusPending,
	}

	mockRepo.On("Create",
		mock.Anything,
		mock.MatchedBy(func(req *repository.CreateBookingRequest) bool {
			return req.ProfessionalID == 1 &&
				req.ClientID == 1 &&
				req.StartTime.Equal(startTime) &&
				req.EndTime.Equal(endTime)
		}),
	).Return(expectedBooking, nil)

	// Execução
	resp, err := svc.CreateBooking(context.Background(), req)

	// Verificações
	assert.NoError(t, err)
	assert.NotNil(t, resp)
	assert.Equal(t, uint(1), resp.ID)

	// Verifica todas as expectativas
	mockRepo.AssertExpectations(t)
}
func TestCreateBooking_PastDate(t *testing.T) {
	mockRepo := new(repository.MockBookingRepository)
	svc := service.NewBookingService(mockRepo)

	req := &service.CreateBookingRequest{
		ProfessionalID: 1,
		ClientID:       1,
		StartTime:      time.Now().Add(-2 * time.Hour),
		EndTime:        time.Now().Add(-1 * time.Hour),
	}

	_, err := svc.CreateBooking(context.Background(), req)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "impossivel agendar no passado")

	// Verifica que os métodos do repositório não foram chamados
	mockRepo.AssertNotCalled(t, "IsTimeSlotAvailable")
	mockRepo.AssertNotCalled(t, "Create")
}

func TestUpdateBookingStatus_Success(t *testing.T) {
	mockRepo := new(repository.MockBookingRepository)
	svc := service.NewBookingService(mockRepo)

	bookingID := uint(1)
	newStatus := domain.StatusConfirmed

	expectedBooking := &domain.Booking{
		ID:            bookingID,
		ProfessionalID: 1,
		ClientID:      1,
		Status:        newStatus,
	}

	mockRepo.On("UpdateStatus", mock.Anything, bookingID, newStatus).
		Return(expectedBooking, nil)

	resp, err := svc.UpdateBookingStatus(context.Background(), bookingID, newStatus)

	assert.NoError(t, err)
	assert.Equal(t, newStatus, resp.Status)
	mockRepo.AssertExpectations(t)
}
