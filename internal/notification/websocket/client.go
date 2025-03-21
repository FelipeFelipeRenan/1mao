package websocket

import (
	"encoding/json"
	"fmt"
	"1mao/internal/notification/domain"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID   int
	Conn *websocket.Conn
	Send chan domain.Message
}

func NewClient(id int, conn *websocket.Conn) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Send: make(chan domain.Message),
	}
}

func (c *Client) Listen() {
	defer func() {
		H.Unregister <- c // Remove cliente ao desconectar
		c.Conn.Close()
	}()

	// Goroutine para receber mensagens e enviá-las ao Hub
	go func() {
		for {
			_, msg, err := c.Conn.ReadMessage()
			if err != nil {
				fmt.Println("❌ Erro ao ler mensagem do WebSocket:", err)
				break
			}

			// Deserializar mensagem JSON
			var message domain.Message
			if err := json.Unmarshal(msg, &message); err != nil {
				fmt.Println("❌ Erro ao desserializar mensagem:", err)
				continue
			}

			message.SenderID = c.ID // Definir ID do remetente
			fmt.Printf("📩 Mensagem recebida de %d para %d: %s\n", message.SenderID, message.ReceiverID, message.Content)

			// Enviar mensagem ao Hub
			H.Broadcast <- message
		}
	}()

	// Loop para enviar mensagens ao cliente
	for msg := range c.Send {
		data, err := json.Marshal(msg)
		if err != nil {
			fmt.Println("❌ Erro ao serializar mensagem:", err)
			continue
		}
		if err := c.Conn.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Println("❌ Erro ao enviar mensagem:", err)
			break
		}
	}
}
