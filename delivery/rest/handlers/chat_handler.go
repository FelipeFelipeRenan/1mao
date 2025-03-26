package handlers

import (
	
	"net/http"
	"strconv"

	"1mao/internal/notification/repository"
	"1mao/internal/notification/websocket"
	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
	"gorm.io/gorm"
)

var chatUpgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func HandleChatWebSocket(w http.ResponseWriter, r *http.Request, db *gorm.DB, hub *websocket.Hub) {
	vars := mux.Vars(r)
	userType := vars["type"]
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	conn, err := chatUpgrader.Upgrade(w, r, nil)
	if err != nil {
		http.Error(w, "Erro ao conectar WebSocket", http.StatusInternalServerError)
		return
	}

	messageRepo := repository.NewMessageRepository(db)
	client := websocket.NewClient(userID, userType, conn, hub, messageRepo)
	hub.Register <- client

	// Enviar mensagens antigas ao conectar
	go func() {
		messages, err := messageRepo.GetMessages(userID, userType, 2, "professional") // Exemplo
		if err == nil {
			for _, msg := range messages {
				client.Send <- msg
			}
		}
	}()

	go client.Listen()
	go client.Write()
}


