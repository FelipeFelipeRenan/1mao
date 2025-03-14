package routes

import (
	"1mao/delivery/rest/handlers"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func HealthRoutes(r *mux.Router, db *gorm.DB) {
	healthHandler := handlers.NewHealthHandler(db)

	r.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")
	r.HandleFunc("/ready", healthHandler.ReadyCheck).Methods("GET")

}
