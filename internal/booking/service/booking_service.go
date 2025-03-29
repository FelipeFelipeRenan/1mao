package service

import (
	"1mao/internal/booking/domain"
	"1mao/internal/booking/repository"
	"context"
	"errors"
	"time"
)

// BookingService define os serviços disponíveis para agendamentos
type BookingService interface {
	CreateBooking(ctx context.Context, req *CreateBookingRequest) (*BookingResponse, error)
	GetBooking(ctx context.Context, id uint) (*BookingResponse, error)
	ListProfessionalBookings(ctx context.Context, professionalID uint, filters *BookingFilters) ([]*BookingResponse, error)
	ListClientBookings(ctx context.Context, clientID uint, filters *BookingFilters) ([]*BookingResponse, error)
		UpdateBookingStatus(ctx context.Context, id uint, status domain.BookingStatus) (*BookingResponse, error)
	CancelBooking(ctx context.Context, id uint) error
}

type bookingService struct {
	bookingRepo repository.BookingRepository
}

func NewBookingService(bookingRepo repository.BookingRepository) BookingService {
	return &bookingService{bookingRepo: bookingRepo}
}

// DTOs
// CreateBookingRequest define o payload para criação de agendamento
// @Model CreateBookingRequest
type CreateBookingRequest struct {
	ProfessionalID uint      `json:"professional_id"`
	ClientID       uint      `json:"client_id"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
}

// BookingResponse define a resposta de agendamento
// @Model BookingResponse
type BookingResponse struct {
	ID             uint                 `json:"id"`
	ProfessionalID uint                 `json:"professional_id"`
	ClientID       uint                 `json:"client_id"`
	StartTime      time.Time            `json:"start_time"`
	EndTime        time.Time            `json:"end_time"`
	Status         domain.BookingStatus `json:"status"`
	CreatedAt      time.Time            `json:"created_at"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

// BookingFilters define o payload para atualização de status
// @Model BookingFilters
type BookingFilters struct {
	From   time.Time
	To     time.Time
	Status domain.BookingStatus
}


func (s *bookingService) CreateBooking(ctx context.Context, req *CreateBookingRequest) (*BookingResponse, error) {

	// validação basica
	if req.StartTime.Before(time.Now()) {
		return nil, errors.New("impossivel agendar no passado")
	}
	if req.EndTime.Before(req.StartTime) {
		return nil, errors.New("o fim do agendamento precisa ser depois do inicio")
	}
	// Verifica disponibilidade
	available, err := s.bookingRepo.IsTimeSlotAvailable(
		ctx,
		req.ProfessionalID,
		req.StartTime,
		req.EndTime,
	)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, domain.ErrTimeSlotUnavailable
	}

	// adicionar no repositorio
	booking, err := s.bookingRepo.Create(ctx, &repository.CreateBookingRequest{
		ProfessionalID: req.ProfessionalID,
		ClientID:       req.ClientID,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return s.toResponse(booking), nil

}

func (s *bookingService) GetBooking(ctx context.Context, id uint) (*BookingResponse, error) {
	bookings, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(bookings), nil
}

func (s *bookingService) ListProfessionalBookings(ctx context.Context, professionalID uint, filters *BookingFilters) ([]*BookingResponse, error) {
	var from, to time.Time
	if filters != nil {
		from = filters.From
		to = filters.To
	}

	bookings, err := s.bookingRepo.ListByProfessional(ctx, professionalID, from, to)
	if err != nil {
		return nil, err // Propagação correta do erro
	}

	var filtered []*domain.Booking
	if filters != nil && filters.Status != "" {
		for _, b := range bookings {
			if b.Status == filters.Status {
				filtered = append(filtered, b)
			}
		}
	} else {
		filtered = bookings
	}

	return s.toListResponse(filtered), nil
}

func (s *bookingService) ListClientBookings(ctx context.Context, clientID uint, filters *BookingFilters) ([]*BookingResponse, error) {
    var from, to time.Time
    if filters != nil {
        from = filters.From
        to = filters.To
    }

    bookings, err := s.bookingRepo.ListByClient(ctx, clientID, from, to)
    if err != nil {
        return nil, err
    }

    var filtered []*domain.Booking
    if filters != nil && filters.Status != "" {
        for _, b := range bookings {
            if b.Status == filters.Status {
                filtered = append(filtered, b)
            }
        }
    } else {
        filtered = bookings
    }

    return s.toListResponse(filtered), nil
}

func (s *bookingService) UpdateBookingStatus(ctx context.Context, id uint, status domain.BookingStatus) (*BookingResponse, error) {
	booking, err := s.bookingRepo.UpdateStatus(ctx, id, status)
	if err != nil {
		return nil, err
	}
	return s.toResponse(booking), nil
}

func (s *bookingService) CancelBooking(ctx context.Context, id uint) error {
	_, err := s.bookingRepo.UpdateStatus(ctx, id, domain.StatusCancelled)
	return err
}

// Helpers
func (s *bookingService) toResponse(booking *domain.Booking) *BookingResponse {
	return &BookingResponse{
		ID:             booking.ID,
		ProfessionalID: booking.ProfessionalID,
		ClientID:       booking.ClientID,
		StartTime:      booking.StartTime,
		EndTime:        booking.EndTime,
		Status:         booking.Status,
		CreatedAt:      booking.CreatedAt,
		UpdatedAt:      booking.UpdatedAt,
	}
}

func (s *bookingService) toListResponse(bookings []*domain.Booking) []*BookingResponse {
	var response []*BookingResponse
	for _, b := range bookings {
		response = append(response, s.toResponse(b))
	}
	return response
}
