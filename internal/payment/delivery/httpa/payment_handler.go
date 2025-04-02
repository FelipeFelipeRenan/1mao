package httpa

import (
	"1mao/internal/payment/dtos"
	"1mao/internal/payment/service"
	"encoding/json"
	"net/http"
)




type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler{
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request){
	var req dtos.CreatePaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil{
		http.Error(w, "corpo da requisição inválido", http.StatusBadRequest)
		return
	}

	transaction, err := h.paymentService.CreatePayment(req.BookingID, req.Amount, req.Method)
	if err != nil {
		http.Error(w, "falha ao criar pagamento", http.StatusInternalServerError)
		return
	}

	response := dtos.PaymentResponse{
		ID: transaction.ID,
		Status: string(transaction.Status),
		GatewayID: transaction.GatewayID,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}