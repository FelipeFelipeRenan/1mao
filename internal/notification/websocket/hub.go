package websocket

import (
	"fmt"
	"sync"
	

	"1mao/internal/notification/domain"
	"1mao/internal/notification/repository"
)

// Hub gerencia as conexões WebSocket
type Hub struct {
	Clients    map[string]*Client
	Broadcast  chan domain.Message
	Register   chan *Client
	Unregister chan *Client
	Repo       *repository.MessageRepository
	mu         sync.Mutex
}

// Criar um novo Hub
func NewHub(repo *repository.MessageRepository) *Hub {
	return &Hub{
		Clients:    make(map[string]*Client),
		Broadcast:  make(chan domain.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Repo:       repo,
	}
}

// Gerenciar conexões e mensagens
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			clientKey := fmt.Sprintf("%s:%d", client.UserType, client.ID)

			h.mu.Lock()
			h.Clients[clientKey] = client
			h.mu.Unlock()

			// Carregar mensagens antigas
			messages, err := h.Repo.GetMessages(client.ID, client.UserType, 2, "professional") // Exemplo: carregando conversas com profissional 2
			if err == nil {
				for _, msg := range messages {
					client.Send <- msg
				}
			}

			fmt.Printf("✅ Cliente %d (%s) registrado\n", client.ID, client.UserType)

		case client := <-h.Unregister:
			clientKey := fmt.Sprintf("%s:%d", client.UserType, client.ID)

			h.mu.Lock()
			if _, ok := h.Clients[clientKey]; ok {
				delete(h.Clients, clientKey)
				close(client.Send)
			}
			h.mu.Unlock()
			fmt.Printf("🚪 Cliente %d desconectado\n", client.ID)

		case msg := <-h.Broadcast:
			receiverKey := fmt.Sprintf("%s:%d", msg.ReceiverType, msg.ReceiverID)

			h.mu.Lock()
			if recipient, ok := h.Clients[receiverKey]; ok {
				recipient.Send <- msg
			} else {
				fmt.Printf("⚠️ Cliente %d (%s) não está online\n", msg.ReceiverID, msg.ReceiverType)
			}

			// Salvar a mensagem no banco
			err := h.Repo.SaveMessage(msg)
			if err != nil {
				fmt.Println("❌ Erro ao salvar mensagem:", err)
			}

			h.mu.Unlock()
		}
	}
}
