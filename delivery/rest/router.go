package rest

import (
	"1mao/delivery/rest/handlers"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)


func NewRouter(db *gorm.DB) http.Handler{
	r := mux.NewRouter()

	healthHandler := handlers.NewHealthHandler(db)
	r.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")
	r.HandleFunc("/ready", healthHandler.ReadyCheck).Methods("GET")

	return r
}