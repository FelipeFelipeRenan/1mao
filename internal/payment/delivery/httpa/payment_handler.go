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
	"github.com/stripe/stripe-go/v81/webhook"
)

type PaymentHandler struct {
	paymentService service.PaymentService
}

func NewPaymentHandler(paymentService service.PaymentService) *PaymentHandler {
	return &PaymentHandler{
		paymentService: paymentService,
	}
}

func (h *PaymentHandler) CreatePayment(w http.ResponseWriter, r *http.Request){
	vars := mux.Vars(r)
	clientID := vars["client_id"]

	var req dtos.CreatePaymentRequest

	if err := json.NewDecoder(r.Body).Decode(&req);err != nil{
		http.Error(w, "requisição invalida", http.StatusBadRequest)
		return
	} 

	if req.Amount <= 0{
		http.Error(w, "o montante deve ser positivo", http.StatusBadRequest)
		return 
	}

	transaction, err := h.paymentService.CreatePayment(clientID, req.BookingID, req.Amount, req.Method)
	if err != nil {
		http.Error(w, "falha ao criar pagamento", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type","application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(transaction)

}

func (h *PaymentHandler) HandleWebhook(w http.ResponseWriter, r *http.Request) {
	log.Println("Webhook recebido")

    const MaxBodyBytes = int64(65536)
    r.Body = http.MaxBytesReader(w, r.Body, MaxBodyBytes)

    payload, err := io.ReadAll(r.Body)
	log.Printf("Payload: %s", string(payload))

    if err != nil {
        http.Error(w, "erro ao ler corpo", http.StatusBadRequest)
        return
    }

	

    signature := r.Header.Get("Stripe-Signature")
    event, err := webhook.ConstructEvent(payload, signature, os.Getenv("STRIPE_WEBHOOK_SECRET"))
    if err != nil {
		log.Printf("Erro na validação: %v", err)

        http.Error(w, "assinatura inválida", http.StatusUnauthorized)
        return
    }

	log.Printf("Tipo de evento: %s", event.Type)

    var handlerErr error // Variável única para tratamento de erros

    switch event.Type {
    case "payment_intent.succeeded":
        var intent stripe.PaymentIntent
        if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
            handlerErr = fmt.Errorf("erro no parsing do evento: %w", err)
            break
        }
        if err := h.paymentService.ConfirmPayment(intent.ID); err != nil {
            handlerErr = fmt.Errorf("falha ao confirmar pagamento: %w", err)
        }

    case "payment_intent.payment_failed":
        var intent stripe.PaymentIntent
        if err := json.Unmarshal(event.Data.Raw, &intent); err != nil {
            handlerErr = fmt.Errorf("erro no parsing do evento: %w", err)
            break
        }
        if err := h.paymentService.FailPayment(intent.ID); err != nil {
            handlerErr = fmt.Errorf("falha ao registrar pagamento falho: %w", err)
        }

    default:
        handlerErr = fmt.Errorf("tipo de evento não tratado: %s", event.Type)
    }

    if handlerErr != nil {
        log.Printf("Erro no webhook: %v", handlerErr)
        http.Error(w, handlerErr.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"processed"}`))
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

func (h *PaymentHandler) GetClientPayments(w http.ResponseWriter, r *http.Request){
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
