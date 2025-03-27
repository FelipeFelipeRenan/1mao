package handlers

import "1mao/internal/booking/service"


type BookingHandler struct {
	bookingService service.BookingService
}

func NewBookingService(bookingService service.BookingService) *BookingHandler{
	return &BookingHandler{bookingService: bookingService}
}


// TODO: implementar metodos do handler
