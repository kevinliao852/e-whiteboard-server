package state

import (
	"fmt"
	"sync"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type RoomState struct {
	mu    sync.RWMutex
	rooms map[string]*core.Room
}

func NewRoomState() *RoomState {
	return &RoomState{
		rooms: make(map[string]*core.Room),
	}
}

func (r *RoomState) ensureRooms() {
	if r.rooms == nil {
		r.rooms = make(map[string]*core.Room)
	}
}

func (r *RoomState) getOrCreateRoom(roomID string) *core.Room {
	r.ensureRooms()

	room, ok := r.rooms[roomID]
	if ok {
		return room
	}

	room = &core.Room{
		ID:           roomID,
		Participants: []core.Participant{},
	}
	r.rooms[roomID] = room
	return room
}

// CreateRoom implements [core.RoomService].
func (r *RoomState) CreateRoom(roomID string) (*core.Room, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	room := r.getOrCreateRoom(roomID)
	return room, nil
}

// JoinRoom implements [core.RoomService].
func (r *RoomState) JoinRoom(roomID string, participant core.Participant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	room := r.getOrCreateRoom(roomID)
	for _, existing := range room.Participants {
		if existing == participant {
			return nil
		}
	}

	room.AddParticipant(participant)
	return nil
}

// LeaveRoom implements [core.RoomService].
func (r *RoomState) LeaveRoom(roomID string, participant core.Participant) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	room, ok := r.rooms[roomID]
	if !ok {
		return fmt.Errorf("room %s not found", roomID)
	}

	room.RemoveParticipant(participant)
	if len(room.Participants) == 0 {
		delete(r.rooms, roomID)
	}
	return nil
}

// BroadcastToRoom implements [core.RoomService].
func (r *RoomState) BroadcastToRoom(roomID string, message string) error {
	r.mu.RLock()
	room, ok := r.rooms[roomID]
	if !ok {
		r.mu.RUnlock()
		return fmt.Errorf("room %s not found", roomID)
	}

	participants := make([]core.Participant, len(room.Participants))
	copy(participants, room.Participants)
	r.mu.RUnlock()

	for _, participant := range participants {
		participant.Notify(message)
	}

	return nil
}

var _ core.RoomService = (*RoomState)(nil)
