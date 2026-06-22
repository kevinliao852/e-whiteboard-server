package state

import (
	"fmt"
	"sort"
	"strings"
	"sync"
	"time"

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
		Name:         humanizeRoomID(roomID),
		Status:       "Waiting for participants",
		LastActivity: time.Now(),
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
	room.Status = buildRoomStatus(len(room.Participants))
	room.LastActivity = time.Now()
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
		return nil
	}

	room.Status = buildRoomStatus(len(room.Participants))
	room.LastActivity = time.Now()
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

	r.mu.Lock()
	if currentRoom, ok := r.rooms[roomID]; ok {
		currentRoom.Status = "Active discussion"
		currentRoom.LastActivity = time.Now()
	}
	r.mu.Unlock()

	return nil
}

// ListRooms implements [core.RoomService].
func (r *RoomState) ListRooms() []core.Room {
	r.mu.RLock()
	defer r.mu.RUnlock()

	rooms := make([]core.Room, 0, len(r.rooms))
	for _, room := range r.rooms {
		roomCopy := core.Room{
			ID:           room.ID,
			Name:         room.Name,
			Status:       room.Status,
			LastActivity: room.LastActivity,
			Participants: make([]core.Participant, len(room.Participants)),
		}
		copy(roomCopy.Participants, room.Participants)
		rooms = append(rooms, roomCopy)
	}

	sort.Slice(rooms, func(i, j int) bool {
		return rooms[i].LastActivity.After(rooms[j].LastActivity)
	})

	return rooms
}

func humanizeRoomID(roomID string) string {
	parts := strings.FieldsFunc(roomID, func(r rune) bool {
		return r == '-' || r == '_' || r == ' '
	})
	if len(parts) == 0 {
		return roomID
	}

	for i, part := range parts {
		if part == "" {
			continue
		}
		parts[i] = strings.ToUpper(part[:1]) + strings.ToLower(part[1:])
	}

	return strings.Join(parts, " ")
}

func buildRoomStatus(participantCount int) string {
	if participantCount <= 1 {
		return "Waiting for participants"
	}

	return "Active collaboration"
}

var _ core.RoomService = (*RoomState)(nil)
