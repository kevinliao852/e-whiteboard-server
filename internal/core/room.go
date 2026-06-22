package core

import "time"

type Room struct {
	ID           string
	Participants []Participant
	Name         string
	Status       string
	LastActivity time.Time
}

type Participant interface {
	Notify(message string)
}

func (r *Room) Broadcast(message string) {
	for _, participant := range r.Participants {
		participant.Notify(message)
	}
}

func (r *Room) AddParticipant(p Participant) {
	r.Participants = append(r.Participants, p)
}

func (r *Room) RemoveParticipant(p Participant) {
	for i, participant := range r.Participants {
		if participant == p {
			r.Participants = append(r.Participants[:i], r.Participants[i+1:]...)
			break
		}
	}
}

type RoomModel interface {
	Create(room *Room) error
	GetByID(roomID string) (*Room, error)
}

type RoomService interface {
	CreateRoom(roomID string) (*Room, error)
	JoinRoom(roomID string, participant Participant) error
	LeaveRoom(roomID string, participant Participant) error
	BroadcastToRoom(roomID string, message string) error
	ListRooms() []Room
}
