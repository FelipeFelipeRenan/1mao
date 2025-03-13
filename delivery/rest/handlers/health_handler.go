package handlers

import (
	"log"
	"net/http"

	"gorm.io/gorm"
)

type HealthHandler struct {
	DB *gorm.DB
}

func NewHealthHandler(db *gorm.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

// HealthCheck verifica a conectividade com o banco de dados
func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	sqlDB, err := h.DB.DB()
	if err != nil {
		log.Println("Erro ao obter conexão do banco de dados:", err)
		http.Error(w, "Banco de dados indisponível", http.StatusInternalServerError)
		return
	}

	if err := sqlDB.Ping(); err != nil {
		log.Println("Banco de dados indisponível:", err)
		http.Error(w, "Banco de dados indisponível", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

// ReadyCheck verifica se a aplicação está pronta para receber tráfego
func (h *HealthHandler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}
