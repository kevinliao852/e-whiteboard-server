package wshub

import (
	"net/http"

	"github.com/cockroachdb/errors"
	"github.com/gorilla/websocket"
)

type Room struct {
	Upgrader    *websocket.Upgrader
	RoomId      string
	Clients     map[*websocket.Conn]bool
	Broadcast   chan []byte
	Register    chan *websocket.Conn
	Unregister  chan *websocket.Conn
}

func NewRoom(roomId string) *Room {
	return &Room{
		Upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		RoomId:     roomId,
		Broadcast:  make(chan []byte),
		Register:   make(chan *websocket.Conn),
		Unregister: make(chan *websocket.Conn),
		Clients:    make(map[*websocket.Conn]bool),
	}
}

func RunRoom(room *Room, errChan *chan error) {
	for {
		select {

		case client := <-room.Register:
			room.Clients[client] = true
		case client := <-room.Unregister:
			if _, ok := room.Clients[client]; ok {
				delete(room.Clients, client)
				client.Close()
			}
		case message := <-room.Broadcast:
			for client := range room.Clients {
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

func StoreMessageToDB(message chan []byte, saver MessageSaver) {
	for {
		select {
		case msg := <-message:
			saver.SaveMessage(msg)
		}
	}
}
