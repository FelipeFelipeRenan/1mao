package routes

import (
	"1mao/internal/payment/delivery/httpa"
	"1mao/internal/payment/service"

	"github.com/gorilla/mux"
)


func PaymentRoutes(r *mux.Router, paymentService *service.PaymentService){
	handler := httpa.NewPaymentHandler(*paymentService)

	r.HandleFunc("/payments", handler.CreatePayment).Methods("POST")
	r.HandleFunc("payments/webhook", handler.HandleWebhook).Methods("POST")
}