package httpa

import (
	"encoding/json"
	"net/http"
	"strconv"

	"1mao/internal/professional/domain"
	"1mao/internal/professional/service"

	"github.com/gorilla/mux"
)

// 🔹 Estrutura do Handler
type ProfessionalHandler struct {
	service service.ProfessionalService
}

// 🔹 Criando um novo Handler
func NewProfessionalHandler(service service.ProfessionalService) *ProfessionalHandler {
	return &ProfessionalHandler{service: service}
}

// 🔹 Registro de Profissional
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

// 🔹 Buscar Profissional por ID
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

// 🔹 Buscar Todos os Profissionais
func (h *ProfessionalHandler) GetAllProfessionals(w http.ResponseWriter, r *http.Request) {
	professionals, err := h.service.GetAllProfessionals()
	if err != nil {
		http.Error(w, "Error fetching professionals", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(professionals)
}

// 🔹 Novo Endpoint: Login para Profissional
func (h *ProfessionalHandler) Login(w http.ResponseWriter, r *http.Request) {
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// 🔹 Chamando o serviço para autenticação
	token, err := h.service.Login(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// 🔹 Respondendo com o token JWT
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
