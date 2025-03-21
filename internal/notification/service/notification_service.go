package service

import (
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
    msg := domain.Message{
        SenderID:   notification.SenderID,
        ReceiverID: notification.ReceiverID,
        Content:    notification.Content,
        Type:       "notification",
    }

    s.hub.Broadcast <- msg
    fmt.Println("🔔 Notificação enviada:", msg)
}
