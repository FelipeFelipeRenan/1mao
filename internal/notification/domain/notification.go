package domain

type Notification struct {
	Type       string `json:"type"`
	ID         int   `json:"id"`
	SenderID   int   `json:"sender_id"`
	ReceiverID int   `json:"receiver_id"`
	Content string `json:"content"`
}
