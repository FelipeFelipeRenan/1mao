package httpa

import (
	"1mao/internal/professional/domain"
	"1mao/internal/professional/service"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type ProfessionalHandler struct {
	service service.ProfessionalService
}

func NewProfessionalHandler(service service.ProfessionalService) *ProfessionalHandler {
	return &ProfessionalHandler{service: service}
}

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


func (h *ProfessionalHandler) GetProfessionalByID(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w,"Invalid ID", http.StatusBadRequest)
		return
	}

	professinal, err := h.service.GetProfessionalByID(uint(id))
	if err != nil {
		http.Error(w, "Professinal not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(professinal)
}

func (h *ProfessionalHandler) GetAllProfessionals(w http.ResponseWriter, r *http.Request){
	professionals, err := h.service.GetAllProfessionals()
	if err != nil {
		http.Error(w, "ERror fetching professional", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(professionals)
}