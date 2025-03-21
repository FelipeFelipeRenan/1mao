package routes

import (
	"1mao/delivery/rest/handlers"
	"1mao/internal/notification/websocket"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// Registrar rotas do Chat
func RegisterChatRoutes(r *mux.Router, db *gorm.DB, hub *websocket.Hub) {
	r.HandleFunc("/ws/chat/client/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleClientChatWebSocket(w, r, db,hub) //  Agora o handler recebe o Hub
	})

	r.HandleFunc("/ws/chat/professional/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleProfessionalChatWebSocket(w, r, db, hub) // Rota para profissionais
	})
}