package routes

import (
	"1mao/delivery/rest/handlers"
	"1mao/internal/notification/websocket"
	"net/http"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

func RegisterChatRoutes(r *mux.Router, db *gorm.DB, hub *websocket.Hub) {
	r.HandleFunc("/ws/chat/{type}/{id}", func(w http.ResponseWriter, r *http.Request) {
		handlers.HandleChatWebSocket(w, r, db, hub)
	})
}
