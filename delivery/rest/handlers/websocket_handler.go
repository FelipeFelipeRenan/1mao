package handlers

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true }, // Permitir conexões de qualquer origem (melhorar isso em prod)
}

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// Atualiza a conexão HTTP para WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println("Erro ao fazer upgrade para WebSocket:", err)
		http.Error(w, "Erro ao estabelecer WebSocket", http.StatusInternalServerError)
		return
	}
	defer conn.Close()

	fmt.Println("🔗 Cliente conectado via WebSocket")

	// Loop para leitura e resposta
	for {
		messageType, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("❌ Erro ao ler mensagem:", err)
			break
		}

		fmt.Printf("📩 Mensagem recebida: %s\n", msg)

		// Responder ao cliente confirmando o recebimento
		response := fmt.Sprintf("✅ Mensagem recebida: %s", msg)
		err = conn.WriteMessage(messageType, []byte(response))
		if err != nil {
			fmt.Println("❌ Erro ao enviar resposta:", err)
			break
		}
	}
}
