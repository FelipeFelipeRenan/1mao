package routes

import (
	"1mao/internal/payment/delivery/httpa"
	"1mao/internal/payment/service"

	"github.com/gorilla/mux"
)

// Rotas parar modulo de pagamentos
func PaymentRoutes(r *mux.Router, paymentService *service.PaymentService) {
	handler := httpa.NewPaymentHandler(*paymentService)

	r.HandleFunc("/payments/webhook", handler.HandleWebhook).Methods("POST")
	r.HandleFunc("/clients/{client_id}/payments", handler.CreatePayment).Methods("POST")
	r.HandleFunc("/payments/{id}", handler.GetPaymentStatus).Methods("GET")
	r.HandleFunc("/clients/{client_id}/payments", handler.GetClientPayments).Methods("GET")
}
