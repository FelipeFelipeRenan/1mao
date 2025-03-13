package handlers

import (
	"database/sql"
	"log"
	"net/http"
)

type HealthHandler struct {
	DB *sql.DB
}

func NewHealthHandler(db *sql.DB) *HealthHandler {
	return &HealthHandler{DB: db}
}

func (h *HealthHandler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	err := h.DB.Ping()
	if err != nil {
		log.Println("Banco de dados indisponivel: ", err)
		http.Error(w, "Banco de dados indisponível", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}

func (h *HealthHandler) ReadyCheck(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("READY"))
}
