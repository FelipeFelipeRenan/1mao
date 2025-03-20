package websocket

import (
	"fmt"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
	hub  *Hub
}

func NewClient(conn *websocket.Conn, hub *Hub) *Client {
	return &Client{
		conn: conn,
		send: make(chan []byte),
		hub:  hub,
	}
}

func (c *Client) ReadMessages(){
	defer func ()  {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("Erro ao ler mensagem", err)
			break
		}
		c.hub.bradcast <- msg
	}
}

func (c *Client) WriteMessages(){
	defer c.conn.Close()
	for msg := range c.send{
		err := c.conn.WriteMessage(websocket.TextMessage, msg)
		if err != nil {
			fmt.Println("Erro ao enviar mensagem")
			break
		}
	}
}

