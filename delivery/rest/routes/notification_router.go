package routes

import (
	"1mao/delivery/rest/notification/handlers"

	"github.com/gorilla/mux"
)

func RegisterNotificationRoutes(r *mux.Router) {
	r.HandleFunc("/ws/notifications", handlers.HandleWebSocket)
}
