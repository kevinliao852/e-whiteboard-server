package state

import (
	"testing"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type fakeParticipant struct {
	messages []string
}

func (p *fakeParticipant) Notify(message string) {
	p.messages = append(p.messages, message)
}

func TestRoomState_CreateJoinBroadcastLeave(t *testing.T) {
	svc := NewRoomState()
	participant := &fakeParticipant{}

	room, err := svc.CreateRoom("room-1")
	if err != nil {
		t.Fatalf("CreateRoom returned error: %v", err)
	}
	if room.ID != "room-1" {
		t.Fatalf("expected room ID room-1, got %q", room.ID)
	}

	if err := svc.JoinRoom("room-1", participant); err != nil {
		t.Fatalf("JoinRoom returned error: %v", err)
	}

	if err := svc.BroadcastToRoom("room-1", "hello"); err != nil {
		t.Fatalf("BroadcastToRoom returned error: %v", err)
	}
	if len(participant.messages) != 1 || participant.messages[0] != "hello" {
		t.Fatalf("expected one broadcast message, got %#v", participant.messages)
	}

	if err := svc.LeaveRoom("room-1", participant); err != nil {
		t.Fatalf("LeaveRoom returned error: %v", err)
	}

	if _, ok := svc.rooms["room-1"]; ok {
		t.Fatalf("expected room to be removed after last participant left")
	}
}

var _ core.Participant = (*fakeParticipant)(nil)
