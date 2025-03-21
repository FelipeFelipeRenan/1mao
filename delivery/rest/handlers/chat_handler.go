package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"1mao/internal/notification/websocket"

	"github.com/gorilla/mux"
	ws "github.com/gorilla/websocket"
)

var chatUpgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Permitir conexões de qualquer origem (melhorar em prod)
}

// HandleChatWebSocket gerencia conexões WebSocket para o chat
func HandleChatWebSocket(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID, err := strconv.Atoi(vars["id"]) // Captura o ID da URL
	if err != nil {
		http.Error(w, "ID inválido", http.StatusBadRequest)
		return
	}

	conn, err := chatUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao fazer upgrade para WebSocket (chat):", err)
		http.Error(w, "Erro ao estabelecer WebSocket", http.StatusInternalServerError)
		return
	}

	fmt.Printf("🔗 Cliente %d conectado ao WebSocket do Chat\n", userID)

	// Criar cliente e registrá-lo no Hub global
	client := websocket.NewClient(userID, conn)
	websocket.H.Register <- client // Registrar cliente no Hub

	defer func() {
		websocket.H.Unregister <- client // Remover cliente ao desconectar
		conn.Close()
	}()

	// Iniciar escuta de mensagens do cliente
	client.Listen()
}
