package ws

import (
	"context"
	"fmt"
	"log"
)

type ParticipantChannel = chan *Participant

// Room represents a chat room that manages WebSocket clients and message broadcasting.
type Room struct {
	RoomId       string
	Participants map[*Participant]bool
	Broadcast    chan []byte
	Register     chan *Participant
	Unregister   chan *Participant
}

// NewRoom creates and returns a new Room instance.
func NewRoom(roomID string) *Room {
	return &Room{
		RoomId:       roomID,
		Broadcast:    make(chan []byte),
		Register:     make(ParticipantChannel),
		Unregister:   make(ParticipantChannel),
		Participants: make(map[*Participant]bool),
	}
}

// Run starts the main loop for the Room to handle client registration, unregistration, and message broadcasting.
func (r *Room) Run(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Shutting down room:", r.RoomId)
			return
		case p := <-r.Register:
			r.Participants[p] = true
		case p := <-r.Unregister:
			if _, ok := r.Participants[p]; ok {
				delete(r.Participants, p)
				if len(r.Participants) == 0 {
					return
				}
				if err := p.Close(); err != nil {
					fmt.Printf("failed to close client: %v\n", err)
				}
			}
		case message := <-r.Broadcast:
			for p := range r.Participants {
				if err := p.WriteMessage(message); err != nil {
					fmt.Printf("failed to write message to client: %v\n", err)
				}
			}
		}
	}
}
