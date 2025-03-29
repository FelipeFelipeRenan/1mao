package handlers

import (
	"1mao/internal/booking/domain"
	"1mao/internal/booking/service"
	"1mao/internal/middleware"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

func NewBookingHandler(bookingService service.BookingService) *BookingHandler {
	return &BookingHandler{bookingService: bookingService, decoder: schema.NewDecoder()}
}

// CreateBookingHandler cria um novo agendamento
// @Summary Cria um novo agendamento
// @Description Cria um novo agendamento entre cliente e profissional
// @Tags Bookings
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param   Authorization   header  string  true  "Token de autenticação (Bearer token)"// @Param booking body service.CreateBookingRequest true "Dados do agendamento"
// @Success 201 {object} service.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /bookings [post]
func (h *BookingHandler) CreateBookingHandler(w http.ResponseWriter, r *http.Request) {
	// Obter informações do usuário autenticado
	claims, ok := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Token inválido")
		return
	}

	userID := uint(claims["user_id"].(float64))
	userRole := claims["role"].(string)

	var req service.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Formato inválido")
		return
	}

	// Validação adicional baseada no perfil
	if userRole == "professional" {
		// Profissional só pode criar bookings para outros clientes
		if req.ProfessionalID != userID {
			respondWithError(w, http.StatusForbidden, "Você só pode criar agendamentos para si mesmo como profissional")
			return
		}
	} else if userRole == "user" {
		// Cliente só pode criar bookings com outros profissionais
		if req.ClientID != userID {
			respondWithError(w, http.StatusForbidden, "Você só pode criar agendamentos para si mesmo como cliente")
			return
		}
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
// @Param   Authorization   header  string  true  "Token de autenticação (Bearer token)"// @Param id path int true "ID do agendamento"
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
// @Param   Authorization   header  string  true  "Token de autenticação (Bearer token)"// @Param from query string false "Data inicial (YYYY-MM-DD)"
// @Param to query string false "Data final (YYYY-MM-DD)"
// @Param status query string false "Status do agendamento"
// @Success 200 {array} service.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /bookings/professional [get]
func (h *BookingHandler) ListProfessionalBookingsHandler(w http.ResponseWriter, r *http.Request) {
	// Obter claims do contexto
	log.Println("Iniciando ListProfessionalBookingsHandler")

	claims, ok := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
	log.Printf("Claims: %+v, ok: %v", claims, ok)
	if !ok {
		respondWithError(w, http.StatusUnauthorized, "Token inválido")
		return
	}

	professionalID := uint(claims["user_id"].(float64))

	// Inicializar filtros
	filters := &service.BookingFilters{}

	// Processar query parameters
	query := r.URL.Query()

	// Filtro de data inicial
	if fromStr := query.Get("from"); fromStr != "" {
		if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
			filters.From = from
		}
	}

	// Filtro de data final
	if toStr := query.Get("to"); toStr != "" {
		if to, err := time.Parse(time.RFC3339, toStr); err == nil {
			filters.To = to
		}
	}

	// Filtro de status
	if statusStr := query.Get("status"); statusStr != "" {
		filters.Status = domain.BookingStatus(statusStr)
	}

	// Chamar service
	bookings, err := h.bookingService.ListProfessionalBookings(r.Context(), professionalID, filters)
	if err != nil {
		handleServiceError(w, err)
		return
	}

	// Retornar array vazio se não houver resultados
	if bookings == nil {
		respondWithJSON(w, http.StatusOK, []interface{}{})
		return
	}

	respondWithJSON(w, http.StatusOK, bookings)
}

// ListClientBookingsHandler lista agendamentos de um cliente
// @Summary Lista agendamentos do cliente
// @Description Retorna a lista de agendamentos de um cliente específico
// @Tags Bookings
// @Security ApiKeyAuth
// @Param   Authorization   header  string  true  "Token de autenticação (Bearer token)"
// @Produce json
// @Param client_id query int true "ID do cliente"
// @Param from query string false "Data inicial (YYYY-MM-DD)"
// @Param to query string false "Data final (YYYY-MM-DD)"
// @Param status query string false "Status do agendamento"
// @Success 200 {array} service.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Router /bookings/client [get]
func (h *BookingHandler) ListClientBookingsHandler(w http.ResponseWriter, r *http.Request) {
    // Obter claims do contexto
    claims, ok := r.Context().Value(middleware.UserContextKey).(jwt.MapClaims)
    if !ok {
		respondWithError(w, http.StatusUnauthorized, "Token inválido")
        return
    }
	
    clientID := uint(claims["user_id"].(float64))
	log.Println("id do cliente: ", clientID)
    
    // Processar query parameters para filtros
    filters := &service.BookingFilters{}
    if fromStr := r.URL.Query().Get("from"); fromStr != "" {
        if from, err := time.Parse(time.RFC3339, fromStr); err == nil {
            filters.From = from
        }
    }
    if toStr := r.URL.Query().Get("to"); toStr != "" {
        if to, err := time.Parse(time.RFC3339, toStr); err == nil {
            filters.To = to
        }
    }
    if statusStr := r.URL.Query().Get("status"); statusStr != "" {
        filters.Status = domain.BookingStatus(statusStr)
    }

    // Chamar service com filtros
    bookings, err := h.bookingService.ListClientBookings(r.Context(), clientID, filters)
    if err != nil {
        handleServiceError(w, err)
        return
    }

    if bookings == nil {
        respondWithJSON(w, http.StatusOK, []interface{}{})
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
// @Param   Authorization   header  string  true  "Token de autenticação (Bearer token)"// @Param id path int true "ID do agendamento"
// @Param status body UpdateStatusRequest true "Novo status"
// @Success 200 {object} service.BookingResponse
// @Failure 400 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Router /bookings/{id}/status [put]
func (h *BookingHandler) UpdateBookingStatusHandler(w http.ResponseWriter, r *http.Request) {
	// DEBUG similar para professional_id
	professionalIDStr := r.URL.Query().Get("professional_id")
	log.Printf("professional_id string: %q", professionalIDStr)

	professionalID, err := strconv.ParseUint(professionalIDStr, 10, 32)
	if err != nil {
		log.Printf("Erro na conversão do professional_id: %v", err)
		respondWithError(w, http.StatusBadRequest, "ID do profissional deve ser um número válido")
		return
	}
	log.Printf("professional_id convertido: %d", professionalID)
	var req struct {
		Status domain.BookingStatus `json:"status"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondWithError(w, http.StatusBadRequest, "Pacote de requisição invalido")
		return
	}
	booking, err := h.bookingService.UpdateBookingStatus(r.Context(), uint(professionalID), req.Status)
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
