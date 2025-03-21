package websocket

import (
	"fmt"
	"sync"

	"1mao/internal/notification/domain"
)

// Hub gerencia as conexões WebSocket
type Hub struct {
	Clients    map[int]*Client  // Agora o índice é um `int` para IDs
	Broadcast  chan domain.Message // Canal de mensagens
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex // Proteção contra concorrência
}

// Criar um novo Hub
func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]*Client),
		Broadcast:  make(chan domain.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

// Método para rodar o Hub e gerenciar clientes e mensagens
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			fmt.Printf("✅ Cliente %d registrado\n", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			fmt.Printf("🚪 Cliente %d desconectado\n", client.ID)

		case msg := <-h.Broadcast:
			fmt.Printf("📢 Broadcast: %d -> %d: %s\n", msg.SenderID, msg.ReceiverID, msg.Content)

			// Encontrar o destinatário e enviar mensagem
			h.mu.Lock()
			if recipient, ok := h.Clients[msg.ReceiverID]; ok {
				recipient.Send <- msg
			} else {
				fmt.Printf("⚠️ Cliente %d não está online\n", msg.ReceiverID)
			}
			h.mu.Unlock()
		}
	}
}
