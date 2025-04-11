package httpa

import (
	"1mao/internal/payment/dtos"
	"1mao/internal/payment/service"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/stripe/stripe-go/v81"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["client_id"]

	var req dtos.CreatePaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "requisição invalida", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "o montante deve ser positivo", http.StatusBadRequest)
		return
	}

	transaction, err := h.paymentService.CreatePayment(clientID, req.BookingID, req.Amount, req.Method)
	if err != nil {
		http.Error(w, "falha ao criar pagamento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)

}

func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	const MaxBodyBytes = int64(65536)

	r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)
	payload, err := io.ReadAll(r.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		return
	}

	event := stripe.Event{}

	if err := json.Unmarshal(payload, &event); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to parse webhook body json: %v\n", err.Error())
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao realizar parsing do webhook: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = h.paymentService.ConfirmPayment(paymentIntent.ID); err != nil {
			log.Println("Erro ao atualizar o status do pagamento: ", err)
		}

		fmt.Println("PaymentIntent was successful!")
	case "payment_method.attached":
		var paymentMethod stripe.PaymentMethod
		err := json.Unmarshal(event.Data.Raw, &paymentMethod)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error parsing webhook JSON: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		fmt.Println("PaymentMethod was attached to a Customer!")
	case "payment_method.failed":
		var paymentIntent stripe.PaymentIntent
		err := json.Unmarshal(event.Data.Raw, &paymentIntent)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Erro ao realizar parsing do webhook: %v\n", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if err = h.paymentService.FailPayment(paymentIntent.ID); err != nil{
			log.Println("Erro ao atualizar o status do pagamento: ", err)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	w.WriteHeader(http.StatusOK)
}

func (h *PaymentHandler) GetPaymentStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	paymentID := vars["id"]

	transaction, err := h.paymentService.GetPaymentByID(paymentID)
	if err != nil {
		http.Error(w, "pagamento nao encontrado", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"status": string(transaction.Status),
		"id":     transaction.ID,
	})
}

func (h *PaymentHandler) GetClientPayments(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID := vars["client_id"]

	payments, err := h.paymentService.GetClientPayments(clientID)
	if err != nil {
		http.Error(w, "falha ao coletar pagamentos", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(payments)

}
