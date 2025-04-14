package routes

import (
	"1mao/delivery/rest/handlers"

	"github.com/gorilla/mux"
)

// Rota para notificação
func RegisterNotificationRoutes(r *mux.Router) {
	r.HandleFunc("/ws/notifications", handlers.HandleNotificationWebSocket)
}
