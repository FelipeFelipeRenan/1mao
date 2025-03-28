package service

import (
	"1mao/internal/booking/domain"
	"1mao/internal/booking/repository"
	"context"
	"errors"
	"time"
)

// BookingService define os serviços disponíveis para agendamentos
// @Service
type BookingService interface {
	CreateBooking(ctx context.Context, req *CreateBookingRequest) (*BookingResponse, error)
	GetBooking(ctx context.Context, id uint) (*BookingResponse, error)
	ListProfessionalBookings(ctx context.Context, professionalID uint, filters *BookingFilters) ([]*BookingResponse, error)
	ListClientBookings(ctx context.Context, clientID uint) ([]*BookingResponse, error)
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
	ServiceID      uint      `json:"service_id"`
	Date           time.Time `json:"date"`
	StartTime      time.Time `json:"start_time"`
	EndTime        time.Time `json:"end_time"`
}

// BookingResponse define a resposta de agendamento
// @Model BookingResponse
type BookingResponse struct {
	ID             uint                 `json:"id"`
	ProfessionalID uint                 `json:"professional_id"`
	ClientID       uint                 `json:"client_id"`
	ServiceID      uint                 `json:"service_id"`
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

// Cria um novo agendamento
// @Summary Cria um novo agendamento
// @Description Cria um novo agendamento entre cliente e profissional
// @Tags Bookings
// @Accept json
// @Produce json
// @Param booking body CreateBookingRequest true "Dados do agendamento"
// @Success 201 {object} BookingResponse "Agendamento criado com sucesso"
// @Failure 400 {object} map[string]string "Dados inválidos"
// @Failure 409 {object} map[string]string "Conflito de horário"
// @Failure 500 {object} map[string]string "Erro interno"
// @Router /bookings [post]
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
		ServiceID:      req.ServiceID,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
	})
	if err != nil {
		return nil, err
	}
	return s.toResponse(booking), nil

}

// Busca um agendamento por ID
// @Summary Busca agendamento por ID
// @Description Retorna os detalhes de um agendamento específico
// @Tags Bookings
// @Produce json
// @Param id path int true "ID do agendamento"
// @Success 200 {object} BookingResponse "Agendamento encontrado"
// @Failure 404 {object} map[string]string "Agendamento não encontrado"
// @Router /bookings/{id} [get]
func (s *bookingService) GetBooking(ctx context.Context, id uint) (*BookingResponse, error) {
	bookings, err := s.bookingRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return s.toResponse(bookings), nil
}

// Lista agendamentos de um profissional
// @Summary Lista agendamentos do profissional
// @Description Retorna a lista de agendamentos de um profissional com filtros opcionais
// @Tags Bookings
// @Produce json
// @Param professional_id query int true "ID do profissional"
// @Param from query string false "Data inicial (YYYY-MM-DD)"
// @Param to query string false "Data final (YYYY-MM-DD)"
// @Param status query string false "Status (pending, confirmed, cancelled, completed)"
// @Success 200 {array} BookingResponse "Lista de agendamentos"
// @Router /bookings/professional [get]
func (s *bookingService) ListProfessionalBookings(ctx context.Context, professionalID uint, filters *BookingFilters) ([]*BookingResponse, error) {
	var from, to time.Time
	var status domain.BookingStatus

	if filters != nil {
		from = filters.From
		to = filters.To
		status = filters.Status
	}

	bookings, err := s.bookingRepo.ListByProfessional(ctx, professionalID, from, to)
	if err != nil {
		return nil, err
	}

	// filtra por status se necessario
	var filtered []*domain.Booking
	if status != "" {
		for _, b := range bookings {
			if b.Status == status {
				filtered = append(filtered, b)
			}

		}
	} else {
		filtered = bookings
	}
	return s.toListResponse(filtered), nil
}

func (s *bookingService) ListClientBookings(ctx context.Context, clientID uint) ([]*BookingResponse, error) {
	bookings, err := s.bookingRepo.ListByClient(ctx, clientID)
	if err != nil {
		return nil, err
	}
	return s.toListResponse(bookings), nil
}

// Atualiza o status de um agendamento
// @Summary Atualiza status do agendamento
// @Description Atualiza o status de um agendamento existente
// @Tags Bookings
// @Accept json
// @Produce json
// @Param id path int true "ID do agendamento"
// @Success 200 {object} BookingResponse "Agendamento atualizado"
// @Failure 400 {object} map[string]string "Transição de status inválida"
// @Router /bookings/{id}/status [put]
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
		ProfessionalID: booking.ProfessinalID,
		ClientID:       booking.ClientID,
		ServiceID:      booking.ServiceID,
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
