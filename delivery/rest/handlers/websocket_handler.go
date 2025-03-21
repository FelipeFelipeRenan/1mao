package handlers

import (
	"fmt"
	"net/http"


	ws "github.com/gorilla/websocket"
)

var notificationUpgrader = ws.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Permitir conexões de qualquer origem (melhorar em prod)
}

// HandleNotificationWebSocket gerencia conexões WebSocket para notificações
func HandleNotificationWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := notificationUpgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao fazer upgrade para WebSocket (notificações):", err)
		http.Error(w, "Erro ao estabelecer WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	fmt.Println("🔗 Cliente conectado ao WebSocket de Notificações")

	// Loop para leitura e resposta
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("❌ Erro ao ler mensagem (notificação):", err)
			break
		}

		fmt.Printf("📩 Notificação recebida: %s\n", msg)

		// Responder ao cliente confirmando o recebimento
		response := fmt.Sprintf("✅ Notificação recebida: %s", msg)
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			fmt.Println("❌ Erro ao enviar resposta (notificação):", err)
			break
		}
	}
}
