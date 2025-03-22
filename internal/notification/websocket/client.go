package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"1mao/internal/notification/domain"
	"1mao/internal/notification/repository"

	"github.com/gorilla/websocket"
)

// Estrutura do cliente WebSocket
type Client struct {
	ID         int
	UserType   string
	Conn       *websocket.Conn
	Send       chan domain.Message
	Hub        *Hub
	Repo       *repository.MessageRepository
}

// Criar um novo cliente WebSocket
func NewClient(id int, userType string, conn *websocket.Conn, hub *Hub, repo *repository.MessageRepository) *Client {
	return &Client{
		ID:       id,
		UserType: userType,
		Conn:     conn,
		Send:     make(chan domain.Message, 256),
		Hub:      hub,
		Repo:     repo,
	}
}

// Método para ouvir mensagens do WebSocket
func (c *Client) Listen() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	for {
		_, msgData, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println("❌ Erro ao ler mensagem:", err)
			break
		}

		var msg domain.Message
		err = json.Unmarshal(msgData, &msg)
		if err != nil {
			log.Println("❌ Erro ao decodificar mensagem:", err)
			continue
		}

		// Definir remetente automaticamente
		msg.SenderID = c.ID
		msg.SenderType = c.UserType
		msg.Timestamp = time.Now().UTC()

		// Enviar mensagem para o hub
		c.Hub.Broadcast <- msg
	}
}

// Método para enviar mensagens para o WebSocket do cliente
func (c *Client) Write() {
	defer c.Conn.Close()

	for msg := range c.Send {
		msgData, err := json.Marshal(msg)
		if err != nil {
			log.Println("❌ Erro ao serializar mensagem:", err)
			continue
		}

		err = c.Conn.WriteMessage(websocket.TextMessage, msgData)
		if err != nil {
			log.Println("❌ Erro ao enviar mensagem:", err)
			break
		}
		fmt.Printf("📤 Mensagem enviada para %d (%s): %s\n", msg.ReceiverID, msg.ReceiverType, msg.Content)
	}
}
