package websocket

import (
	"fmt"
	"1mao/internal/notification/domain"
	"sync"
)

type Hub struct {
	Clients    map[int]*Client
	Broadcast  chan domain.Message
	Register   chan *Client
	Unregister chan *Client
	mu         sync.Mutex
}

var H = NewHub()

func NewHub() *Hub {
	return &Hub{
		Clients:    make(map[int]*Client),
		Broadcast:  make(chan domain.Message),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mu.Lock()
			h.Clients[client.ID] = client
			h.mu.Unlock()
			fmt.Printf("✅ Cliente %d registrado no Hub\n", client.ID)

		case client := <-h.Unregister:
			h.mu.Lock()
			if _, ok := h.Clients[client.ID]; ok {
				delete(h.Clients, client.ID)
				close(client.Send)
			}
			h.mu.Unlock()
			fmt.Printf("⛔ Cliente %d desconectado\n", client.ID)

		case message := <-h.Broadcast:
			h.mu.Lock()
			receiver, exists := h.Clients[message.ReceiverID]
			h.mu.Unlock()

			if exists {
				fmt.Printf("📤 Enviando mensagem de %d para %d: %s\n", message.SenderID, message.ReceiverID, message.Content)
				receiver.Send <- message
			} else {
				fmt.Printf("⚠️ Cliente %d não está conectado, mensagem não enviada\n", message.ReceiverID)
			}
		}
	}
}
