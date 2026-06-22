package core

type ChatMessage struct {
	ID         int    `json:"id"`
	RoomID     string `json:"room-id"`
	SenderID   int    `json:"sender-id"`
	SenderName string `json:"sender-name"`
	Message    string `json:"message"`
}

type ChatService interface {
	ListMessages(roomID string) []ChatMessage
	AppendMessage(roomID string, senderID int, senderName string, message string) ChatMessage
}
