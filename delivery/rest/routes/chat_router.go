package routes

import (
	"1mao/delivery/rest/handlers"

	"github.com/gorilla/mux"
)

// RegisterChatRoutes define a rota do WebSocket para o chat
func RegisterChatRoutes(r *mux.Router) {
	r.HandleFunc("/ws/chat/{id:[0-9]+}", handlers.HandleChatWebSocket)
}
