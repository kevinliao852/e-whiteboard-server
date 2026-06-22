package controllers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type roomParticipantStub struct{}

func (p *roomParticipantStub) Notify(message string) {}

type roomServiceStub struct {
	rooms []core.Room
}

func (s *roomServiceStub) CreateRoom(roomID string) (*core.Room, error) { return nil, nil }
func (s *roomServiceStub) JoinRoom(roomID string, participant core.Participant) error {
	return nil
}
func (s *roomServiceStub) LeaveRoom(roomID string, participant core.Participant) error {
	return nil
}
func (s *roomServiceStub) BroadcastToRoom(roomID string, message string) error { return nil }
func (s *roomServiceStub) ListRooms() []core.Room                              { return s.rooms }

func TestRoomController_ListRooms(t *testing.T) {
	gin.SetMode(gin.TestMode)

	ctrl := NewRoomController(&roomServiceStub{
		rooms: []core.Room{
			{
				ID:           "design-review",
				Name:         "Design Review",
				Status:       "Commenting on navigation states",
				LastActivity: time.Now().Add(-2 * time.Minute),
				Participants: []core.Participant{
					&roomParticipantStub{},
					&roomParticipantStub{},
					&roomParticipantStub{},
					&roomParticipantStub{},
					&roomParticipantStub{},
				},
			},
		},
	})

	router := gin.Default()
	router.GET("/rooms", ctrl.ListRooms)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/rooms", nil)
	router.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var response []RoomResponse
	if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}

	if len(response) != 1 {
		t.Fatalf("expected 1 room, got %d", len(response))
	}

	room := response[0]
	if room.ID != "design-review" {
		t.Fatalf("expected room id design-review, got %q", room.ID)
	}
	if room.Name != "Design Review" {
		t.Fatalf("expected room name Design Review, got %q", room.Name)
	}
	if room.Status != "Commenting on navigation states" {
		t.Fatalf("expected room status to be preserved, got %q", room.Status)
	}
	if room.Participants != 5 {
		t.Fatalf("expected 5 participants, got %d", room.Participants)
	}
	if room.Activity != "2m ago" {
		t.Fatalf("expected activity 2m ago, got %q", room.Activity)
	}
}
