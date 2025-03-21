package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"1mao/internal/notification/websocket"

	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
)

// Configuração do WebSocket
var chatUpgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Manipulador de WebSocket para Chat
func HandleChatWebSocket(w http.ResponseWriter, r *http.Request, hub *websocket.Hub) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	conn, err := chatUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("❌ Erro ao conectar WebSocket:", err)
		http.Error(w, "Erro ao conectar WebSocket", http.StatusInternalServerError)
		return
	}

	// Criar um novo cliente e registrá-lo no Hub
	client := websocket.NewClient(userID, conn, hub)
	hub.Register <- client

	fmt.Printf("🔗 Cliente %d conectado ao WebSocket do Chat\n", userID)

	// Inicia as goroutines para leitura e escrita
	go client.Listen()
	go client.Write()
}
