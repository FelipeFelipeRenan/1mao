package routes

import (
	"1mao/delivery/rest/handlers"
	"1mao/internal/booking/service"

	"github.com/gorilla/mux"
)


func BookingRoutes(r *mux.Router, bookingService service.BookingService){
	handler := handlers.NewBookingService(bookingService)

	// Rotas de Booking
    r.HandleFunc("/bookings", handler.CreateBookingHandler).Methods("POST")
    r.HandleFunc("/bookings/{id}", handler.GetBookingHandler).Methods("GET")
    r.HandleFunc("/bookings/professional", handler.ListProfessionalBookingsHandler).Methods("GET")
    r.HandleFunc("/bookings/client", handler.ListClientBookingsHandler).Methods("GET")
    r.HandleFunc("/bookings/{id}/status", handler.UpdateBookingStatusHandler).Methods("PUT")}