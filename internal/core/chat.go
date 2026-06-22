package core

type ChatMessage struct {
	ID      int    `json:"id"`
	RoomID  string `json:"room-id"`
	Message string `json:"message"`
}

type ChatService interface {
	ListMessages(roomID string) []ChatMessage
	AppendMessage(roomID string, message string) ChatMessage
}
