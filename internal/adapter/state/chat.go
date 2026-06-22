package state

import (
	"sync"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type ChatState struct {
	mu       sync.RWMutex
	messages map[string][]core.ChatMessage
	nextID   int
}

func NewChatState() *ChatState {
	return &ChatState{
		messages: make(map[string][]core.ChatMessage),
		nextID:   1,
	}
}

func (s *ChatState) ListMessages(roomID string) []core.ChatMessage {
	s.mu.RLock()
	defer s.mu.RUnlock()

	history := s.messages[roomID]
	result := make([]core.ChatMessage, len(history))
	copy(result, history)
	return result
}

func (s *ChatState) AppendMessage(roomID string, senderID int, senderName string, message string) core.ChatMessage {
	s.mu.Lock()
	defer s.mu.Unlock()

	chatMessage := core.ChatMessage{
		ID:         s.nextID,
		RoomID:     roomID,
		SenderID:   senderID,
		SenderName: senderName,
		Message:    message,
	}
	s.nextID++

	s.messages[roomID] = append(s.messages[roomID], chatMessage)
	return chatMessage
}

var _ core.ChatService = (*ChatState)(nil)
