package core

type Room struct {
	ID           int
	Participants []Participant
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
	GetByID(roomID int) (*Room, error)
}

type RoomWSModel interface {
	AddToRoom(roomID int, participant Participant) error
	RemoveFromRoom(roomID int, participant Participant) error
	BroadcastToRoom(roomID int, message string) error
}

type RoomService interface {
	CreateRoom() (*Room, error)
	JoinRoom(roomID int, participant Participant) error
	LeaveRoom(roomID int, participant Participant) error
	BroadcastToRoom(roomID int, message string) error
}
