package routes

import (
	"1mao/delivery/rest/handlers"
	"1mao/internal/booking/service"
	"1mao/internal/middleware"

	"github.com/gorilla/mux"
)

// Rotas para Agendamentos
func BookingRoutes(r *mux.Router, bookingService service.BookingService) {
    handler := handlers.NewBookingHandler(bookingService)

    // Rotas para profissionais
    professionalRouter := r.PathPrefix("/professional").Subrouter()
    professionalRouter.Use(middleware.AuthMiddleware("professional"))
    
    professionalRouter.HandleFunc("/bookings/all", handler.ListProfessionalBookingsHandler).Methods("GET")
    professionalRouter.HandleFunc("/bookings/{id:[0-9]+}", handler.GetBookingHandler).Methods("GET")
    professionalRouter.HandleFunc("/bookings/{id:[0-9]+}/status", handler.UpdateBookingStatusHandler).Methods("PUT")

    // Rotas para clientes
    clientRouter := r.PathPrefix("/client").Subrouter()
    clientRouter.Use(middleware.AuthMiddleware("user"))
    clientRouter.HandleFunc("/bookings/all", handler.ListClientBookingsHandler).Methods("GET")

    // Rota compartilhada para criação
    authRouter := r.PathPrefix("").Subrouter()
    authRouter.Use(middleware.AuthMiddleware("user", "professional"))
    authRouter.HandleFunc("/bookings", handler.CreateBookingHandler).Methods("POST")
}