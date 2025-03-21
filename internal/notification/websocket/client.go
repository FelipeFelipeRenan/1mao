package websocket

import (
	
	"fmt"
	"github.com/gorilla/websocket"
	"1mao/internal/notification/domain"
)

// Estrutura do cliente WebSocket
type Client struct {
	ID   int
	Conn *websocket.Conn
	Send chan domain.Message // Canal para envio de mensagens
	Hub  *Hub
}

// Criar um novo cliente WebSocket
func NewClient(id int, conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		ID:   id,
		Conn: conn,
		Send: make(chan domain.Message),
		Hub:  hub,
	}
}

// Método para escutar mensagens enviadas pelo cliente e repassá-las ao Hub
func (c *Client) Listen() {
	defer func() {
		c.Hub.Unregister <- c // Remove o cliente do Hub ao desconectar
		c.Conn.Close()
	}()

	for {
		var msg domain.Message
		err := c.Conn.ReadJSON(&msg) // Lendo JSON do cliente
		if err != nil {
			fmt.Println("❌ Erro ao ler mensagem do WebSocket:", err)
			break
		}

		// Define o ID do remetente como o próprio cliente
		msg.SenderID = c.ID
		fmt.Printf("📩 Mensagem recebida de %d: %+v\n", c.ID, msg)

		// Enviar mensagem ao Hub para distribuição
		c.Hub.Broadcast <- msg
	}
}

// Método para enviar mensagens ao cliente
func (c *Client) Write() {
	defer c.Conn.Close()
	for msg := range c.Send {
		err := c.Conn.WriteJSON(msg)
		if err != nil {
			fmt.Println("❌ Erro ao enviar mensagem via WebSocket:", err)
			break
		}
	}
}
