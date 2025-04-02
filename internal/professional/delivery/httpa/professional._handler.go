package httpa

import (
	"encoding/json"
	"net/http"
	"strconv"

	"1mao/internal/professional/domain"
	"1mao/internal/professional/service"

	"github.com/gorilla/mux"
)

// RegisterRequest define a estrutura para registro de clientes
//
//	@Description	Dados necessários para registrar um novo cliente
type RegisterRequest struct {
	Name     string `json:"name" example:"João Silva"`
	Email    string `json:"email" example:"cliente@example.com"`
	Password string `json:"password" example:"senhaSegura123"`
	Phone    string `json:"phone" example:"+5511999999999"`
}

// LoginRequest define a estrutura para login de clientes
//
//	@Description	Credenciais para autenticação do cliente
type LoginRequest struct {
	Email    string `json:"email" example:"cliente@example.com"`
	Password string `json:"password" example:"senhaSegura123"`
}

// ProfessionalHandler lida com requisições relacionadas a profissionais
//
//	@title			Professional API
//	@version		1.0
//	@description	API para gestão de profissionais
type ProfessionalHandler struct {
	service service.ProfessionalService
}

// NewProfessionalHandler cria uma nova instância do handler
func NewProfessionalHandler(service service.ProfessionalService) *ProfessionalHandler {
	return &ProfessionalHandler{service: service}
}

// Register godoc
//
//	@Summary		Registrar novo profissional
//	@Description	Cria uma nova conta de profissional
//	@Tags			Professionals
//	@Accept			json
//	@Produce		json
//	@Param			professional	body		RegisterRequest	true	"Dados do profissional"
//	@Success		201				{object}	domain.Professional
//	@Failure		400				{object}	map[string]string	"Dados inválidos"
//	@Failure		500				{object}	map[string]string	"Erro interno"
//	@Router			/professional/register [post]
func (h *ProfessionalHandler) Register(w http.ResponseWriter, r *http.Request) {
	var professional domain.Professional
	if err := json.NewDecoder(r.Body).Decode(&professional); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := h.service.Register(&professional); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(professional)
}

// GetProfessionalByID godoc
//
//	@Summary		Obter profissional por ID
//	@Description	Retorna os detalhes de um profissional específico
//	@Tags			Professionals
//	@Produce		json
//	@Param			id	path		int	true	"ID do profissional"
//	@Success		200	{object}	domain.Professional
//	@Failure		400	{object}	map[string]string	"ID inválido"
//	@Failure		404	{object}	map[string]string	"Profissional não encontrado"
//	@Router			/professional/{id} [get]
func (h *ProfessionalHandler) GetProfessionalByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid ID", http.StatusBadRequest)
		return
	}

	professional, err := h.service.GetProfessionalByID(uint(id))
	if err != nil {
		http.Error(w, "Professional not found", http.StatusNotFound)
		return
	}

	json.NewEncoder(w).Encode(professional)
}

// GetAllProfessionals godoc
//
//	@Summary		Listar todos os profissionais
//	@Description	Retorna uma lista com todos os profissionais cadastrados
//	@Tags			Professionals
//	@Produce		json
//	@Success		200	{array}		domain.Professional
//	@Failure		500	{object}	map[string]string	"Erro interno"
//	@Router			/professionals [get]
func (h *ProfessionalHandler) GetAllProfessionals(w http.ResponseWriter, r *http.Request) {
	professionals, err := h.service.GetAllProfessionals()
	if err != nil {
		http.Error(w, "Error fetching professionals", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(professionals)
}

// Login godoc
//
//	@Summary		Login de profissional
//	@Description	Autentica um profissional e retorna um token JWT
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		LoginRequest	true	"Credenciais de login"
//	@Success		200		{object}	map[string]string
//	@Failure		400		{object}	map[string]string
//	@Failure		401		{object}	map[string]string
//	@Router			/professional/login [post]
func (h *ProfessionalHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	token, err := h.service.Login(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
