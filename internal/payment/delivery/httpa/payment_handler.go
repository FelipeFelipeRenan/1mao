package httpa

import (
	"1mao/internal/payment/dtos"
	"1mao/internal/payment/service"
	"encoding/json"
	"io"
	"net/http"
	"os"

	"github.com/stripe/stripe-go/v81"
	"github.com/stripe/stripe-go/v81/webhook"
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


func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request){
	const MaxBodyBytes = int64(65536)
	r.Body = http.MaxBytesReader(w,r.Body, MaxBodyBytes)


	payload, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w,"erro ao ler corpo", http.StatusBadRequest)
		return
	}

	signature := r.Header.Get("Stripe-Signature")
	event, err := webhook.ConstructEvent(payload, signature, os.Getenv("STRIPE_WEBHOOK_SECRET"))
	if err != nil {
		http.Error(w,"assinatura inválida", http.StatusUnauthorized)
		return
	}

	switch event.Type {
	case "payment_intent.succeeded":
		var intent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &intent)
		if err != nil {
			http.Error(w, "erro no parsing do evento", http.StatusBadRequest)
			return
		}
		err = h.paymentService.ConfirmPayment(intent.ID)

	case "payment_intent.payment_failed":
		var intent stripe.PaymentIntent
		json.Unmarshal(event.Data.Raw, &intent)
		err = h.paymentService.FailPayment(intent.ID)

	}
	if err != nil {
		http.Error(w, "falha ao processar evento", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"processed"}`))
}