package service

import (
	"encoding/json"
	"fmt"
	"1mao/internal/notification/domain"
	"1mao/internal/notification/websocket"
)

type NotificationService struct {
	hub *websocket.Hub
}

func NewNotificationService(hub *websocket.Hub) *NotificationService {
	return &NotificationService{hub: hub}
}

func (s *NotificationService) SendNotification(notification domain.Notification) {
	msg, err := json.Marshal(notification)
	if err != nil {
		fmt.Println("Erro ao serializar notificação:", err)
		return
	}
	s.hub.Broadcast <- msg

}
