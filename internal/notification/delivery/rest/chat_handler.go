package rest

import (
	"encoding/json"
	"net/http"
	"strconv"

	_ "1mao/internal/notification/domain"
	"1mao/internal/notification/repository"
)

// ChatHandler gerencia o histórico de mensagens
type ChatHandler struct {
	MessageRepo *repository.MessageRepository
}

// NewChatHandler cria um novo handler para o chat
func NewChatHandler(repo *repository.MessageRepository) *ChatHandler {
	return &ChatHandler{MessageRepo: repo}
}

// @Summary Buscar mensagens de chat
// @Description Retorna o histórico de mensagens entre usuários
// @Tags Chat
// @Produce json
// @Param sender_id query int true "ID do remetente" example(1)
// @Param sender_type query string true "Tipo do remetente" Enums(client, professional)
// @Param receiver_id query int true "ID do destinatário" example(2)
// @Param receiver_type query string true "Tipo do destinatário" Enums(client, professional)
// @Success 200 {array} domain.Message
// @Router /chat/messages [get]
func (h *ChatHandler) GetChatMessages(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	senderID, _ := strconv.Atoi(query.Get("sender_id"))
	senderType := query.Get("sender_type")
	receiverID, _ := strconv.Atoi(query.Get("receiver_id"))
	receiverType := query.Get("receiver_type")

	if senderID == 0 || receiverID == 0 || senderType == "" || receiverType == "" {
		http.Error(w, "Parâmetros inválidos", http.StatusBadRequest)
		return
	}

	messages, err := h.MessageRepo.GetMessages((senderID), senderType, (receiverID), receiverType)
	if err != nil {
		http.Error(w, "Erro ao buscar mensagens", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(messages)
}
