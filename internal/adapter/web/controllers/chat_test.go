package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type chatRoomServiceStub struct{}

func (s *chatRoomServiceStub) CreateRoom(roomID string) (*core.Room, error) { return nil, nil }
func (s *chatRoomServiceStub) JoinRoom(roomID string, participant core.Participant) error {
	return nil
}
func (s *chatRoomServiceStub) LeaveRoom(roomID string, participant core.Participant) error {
	return nil
}
func (s *chatRoomServiceStub) BroadcastToRoom(roomID string, message string) error { return nil }
func (s *chatRoomServiceStub) ListRooms() []core.Room                              { return nil }

type chatServiceStub struct {
	messages map[string][]core.ChatMessage
}

func (s *chatServiceStub) ListMessages(roomID string) []core.ChatMessage {
	return s.messages[roomID]
}

func (s *chatServiceStub) AppendMessage(roomID string, message string) core.ChatMessage {
	return core.ChatMessage{}
}

func TestChatController_GetChatMessages(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := NewChatController(&chatRoomServiceStub{}, &chatServiceStub{
		messages: map[string][]core.ChatMessage{
			"design-review": {
				{
					ID:      1,
					RoomID:  "design-review",
					Message: "Let's review the navigation flow before we finalize the layout.",
				},
				{
					ID:      2,
					RoomID:  "design-review",
					Message: "I'll mark the confusing steps directly on the board.",
				},
			},
		},
	})

	router := gin.Default()
	router.GET("/chat-messages", ctrl.GetChatMessages)

	t.Run("returns history", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/chat-messages?room-id=design-review", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var response []core.ChatMessage
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(response) != 2 {
			t.Fatalf("expected 2 messages, got %d", len(response))
		}
		if response[0].RoomID != "design-review" {
			t.Fatalf("expected room id design-review, got %q", response[0].RoomID)
		}
		if response[0].Message == "" {
			t.Fatal("expected first message to be populated")
		}
	})

	t.Run("returns empty history", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/chat-messages?room-id=empty-room", nil)
		router.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var response []core.ChatMessage
		if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
			t.Fatalf("failed to decode response: %v", err)
		}

		if len(response) != 0 {
			t.Fatalf("expected empty history, got %d messages", len(response))
		}
	})
}
