package handlers

import (
	"1mao/internal/booking/domain"
	"1mao/internal/booking/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/gorilla/schema"
)

// @Model ErrorResponse
type ErrorResponse struct {
	// Mensagem de erro
	// @Example Resource not found
	Error string `json:"error"`
}

// @Model UpdateStatusRequest
type UpdateStatusRequest struct {
	// Novo status do agendamento
	// @Enum pending,confirmed,cancelled,completed
	// @Example confirmed
	Status string `json:"status"`
}

type BookingHandler struct {
	bookingService service.BookingService
	decoder        *schema.Decoder
}

func NewBookingService(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService, decoder: schema.NewDecoder()}
}

// CreateBookingHandler cria um novo agendamento
// @Summary Cria um novo agendamento
// @Description Cria um novo agendamento entre cliente e profissional
// @Tags Bookings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param booking body service.CreateBookingRequest true "Dados do agendamento"
// @Success 201 {object} service.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings [post]
func (h *BookingHandler) CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
    var req service.CreateBookingRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        respondWithError(w, http.StatusBadRequest, "Invalid request payload")
        return
    }

    // Validação básica
    if req.ProfessionalID == 0 || req.ClientID == 0 {
        respondWithError(w, http.StatusBadRequest, "professional_id and client_id are required")
        return
    }

    booking, err := h.bookingService.CreateBooking(r.Context(), &req)
    if err != nil {
        handleServiceError(w, err)
        return
    }

    respondWithJSON(w, http.StatusCreated, booking)
}

// GetBookingHandler obtém um agendamento por ID
// @Summary Obtém um agendamento
// @Description Retorna os detalhes de um agendamento específico
// @Tags Bookings
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do agendamento"
// @Success 200 {object} service.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /bookings/{id} [get]
func (h *BookingHandler) GetBookingHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de agendamento invalido")
		return
	}
	booking, err := h.bookingService.GetBooking(r.Context(), uint(id))
	if err != nil {
		handleServiceError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, booking)
}

// ListProfessionalBookingsHandler lista agendamentos de um profissional
// @Summary Lista agendamentos do profissional
// @Description Retorna a lista de agendamentos de um profissional
// @Tags Bookings
// @Produce json
// @Security ApiKeyAuth
// @Param from query string false "Data inicial (YYYY-MM-DD)"
// @Param to query string false "Data final (YYYY-MM-DD)"
// @Param status query string false "Status do agendamento"
// @Success 200 {array} service.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /bookings/professional [get]
func (h *BookingHandler) ListProfessionalBookingsHandler(w http.ResponseWriter, r *http.Request) {
    professionalIDStr := r.URL.Query().Get("professional_id")
    if professionalIDStr == "" {
        respondWithError(w, http.StatusBadRequest, "O parâmetro professional_id é obrigatório")
        return
    }

    professionalID, err := strconv.ParseUint(professionalIDStr, 10, 32)
    if err != nil {
        respondWithError(w, http.StatusBadRequest, "ID do profissional inválido")
        return
    }
	if err := r.ParseForm(); err != nil {
		respondWithError(w, http.StatusBadRequest, "parametros de busca invalidos")
		return
	}
	var filters service.BookingFilters
	if err := h.decoder.Decode(&filters, r.Form); err != nil {
		respondWithError(w, http.StatusBadRequest, "parametros de filtro invalidos")
		return
	}

	bookings, err := h.bookingService.ListProfessionalBookings(r.Context(), uint(professionalID), &filters)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, bookings)
}

// ListClientBookingsHandler lista agendamentos de um cliente
// @Summary Lista agendamentos do cliente
// @Description Retorna a lista de agendamentos de um cliente específico
// @Tags Bookings
// @Produce json
// @Param client_id query int true "ID do cliente"
// @Param from query string false "Data inicial (YYYY-MM-DD)"
// @Param to query string false "Data final (YYYY-MM-DD)"
// @Param status query string false "Status do agendamento"
// @Success 200 {array} service.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Router /bookings/client [get]
func (h *BookingHandler) ListClientBookingsHandler(w http.ResponseWriter, r *http.Request) {
	clientIDStr := r.URL.Query().Get("client_id")
	if clientIDStr == "" {
		respondWithError(w, http.StatusBadRequest, "O parâmetro client_id é obrigatório")
		return
	}

	clientID, err := strconv.ParseUint(clientIDStr, 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID do cliente inválido")
		return
	}

	var filters service.BookingFilters
	if err := h.decoder.Decode(&filters, r.Form); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid filter parameters")
		return
	}

	bookings, err := h.bookingService.ListClientBookings(r.Context(), uint(clientID))
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Filtra por status se fornecido
	if status := r.URL.Query().Get("status"); status != "" {
		var filtered []*service.BookingResponse
		for _, b := range bookings {
			if string(b.Status) == status {
				filtered = append(filtered, b)
			}
		}
		respondWithJSON(w, http.StatusOK, filtered)
		return
	}

	respondWithJSON(w, http.StatusOK, bookings)
}

// UpdateBookingStatusHandler atualiza o status de um agendamento
// @Summary Atualiza status do agendamento
// @Description Atualiza o status de um agendamento existente
// @Tags Bookings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path int true "ID do agendamento"
// @Param status body UpdateStatusRequest true "Novo status"
// @Success 200 {object} service.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /bookings/{id}/status [put]
func (h *BookingHandler) UpdateBookingStatusHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "ID de agendamento inválido")
		return
	}

	var req struct {
		Status domain.BookingStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Pacote de requisição invalido")
		return
	}
	booking, err := h.bookingService.UpdateBookingStatus(r.Context(), uint(id), req.Status)
	if err != nil {
		handleServiceError(w, err)
		return
	}
	respondWithJSON(w, http.StatusOK, booking)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func handleServiceError(w http.ResponseWriter, err error) {
	switch err {
	case domain.ErrBookingNotFound:
		respondWithError(w, http.StatusNotFound, "Booking not found")
	case domain.ErrTimeSlotUnavailable:
		respondWithError(w, http.StatusConflict, "Time slot unavailable")
	case domain.ErrInvalidStatusTransition:
		respondWithError(w, http.StatusBadRequest, "Invalid status transition")
	default:
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}
}

// TODO: implementar metodos do handler
