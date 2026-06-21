package service

import "github.com/kevinliao852/e-whiteboard-server/internal/core"

type RoomSVC struct {
	RoomModel   core.RoomModel
	RoomWSModel core.RoomWSModel
	Whiteboard  core.WhiteboardModel
}

// BroadcastToRoom implements [core.RoomService].
func (r *RoomSVC) BroadcastToRoom(roomID int, message string) error {
	panic("unimplemented")
}

// CreateRoom implements [core.RoomService].
func (r *RoomSVC) CreateRoom() (*core.Room, error) {
	room := &core.Room{
		Participants: []core.Participant{},
	}

	return room, nil
}

// JoinRoom implements [core.RoomService].
func (r *RoomSVC) JoinRoom(roomID int, participant core.Participant) error {
	panic("unimplemented")
}

// LeaveRoom implements [core.RoomService].
func (r *RoomSVC) LeaveRoom(roomID int, participant core.Participant) error {
	panic("unimplemented")
}

func NewRoomSVC(roomModel core.RoomModel, roomWSModel core.RoomWSModel, whiteboard core.WhiteboardModel) *RoomSVC {
	return &RoomSVC{
		RoomModel:   roomModel,
		RoomWSModel: roomWSModel,
		Whiteboard:  whiteboard,
	}
}

var _ core.RoomService = (*RoomSVC)(nil)
