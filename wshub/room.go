package wshub

import (
	"fmt"
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gorilla/websocket"
)

type RoomActionChannel = chan *websocket.Conn

// Room represents a chat room that manages WebSocket clients and message broadcasting.
type Room struct {
	Upgrader   *websocket.Upgrader
	RoomId     string
	Clients    map[*websocket.Conn]bool
	Broadcast  chan []byte
	Register   chan *websocket.Conn
	Unregister chan *websocket.Conn
}

// NewRoom creates and returns a new Room instance.
func NewRoom(roomID string) *Room {
	return &Room{
		Upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		RoomId:     roomID,
		Broadcast:  make(chan []byte),
		Register:   make(RoomActionChannel),
		Unregister: make(RoomActionChannel),
		Clients:    make(map[*websocket.Conn]bool),
	}
}

// Run starts the main loop for the Room to handle client registration, unregistration, and message broadcasting.
func (r *Room) Run(errChan *chan error) {
	for {
		select {
		case client := <-r.Register:
			r.Clients[client] = true
		case client := <-r.Unregister:
			if _, ok := r.Clients[client]; ok {
				delete(r.Clients, client)
				if err := client.Close(); err != nil {
					fmt.Printf("failed to close client: %v\n", err)
				}
			}
		case message := <-r.Broadcast:
			for client := range r.Clients {
				err := client.WriteMessage(websocket.TextMessage, message)

				if err != nil {
					*errChan <- errors.Wrap(err, "Error writing message")
				}
			}
		}
	}
}

type MessageSaver interface {
	SaveMessage(message []byte) error
}

func StoreMessageToDB(message chan []byte, saver MessageSaver, errChan *chan error) {
	done := make(chan struct{})
	// TODO: graceful shutdown

	for {
		select {
		case msg := <-message:
			err := saver.SaveMessage(msg)

			if err != nil {
				*errChan <- errors.Wrap(err, "Error saving message")
			}
		case <-done:
			return
		}
	}

}
