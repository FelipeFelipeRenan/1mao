package domain

// Message representa uma mensagem de chat entre usuarios
type Message struct {
	SenderID   int    `json:"sender_id"`
	ReceiverID int    `json:"receiver_id"`
	Content    string `json:"content"`
	Type       string `json:"type"`
}
