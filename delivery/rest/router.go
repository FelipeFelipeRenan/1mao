package rest

import (
	"1mao/delivery/rest/handlers"
	"database/sql"
	"net/http"

	"github.com/gorilla/mux"
)


func NewRouter(db *sql.DB) http.Handler{
	r := mux.NewRouter()

	healthHandler := handlers.NewHealthHandler(db)
	r.HandleFunc("/health", healthHandler.HealthCheck).Methods("GET")
	r.HandleFunc("/ready", healthHandler.ReadyCheck).Methods("GET")

	return r
}