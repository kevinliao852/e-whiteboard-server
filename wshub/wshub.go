package wshub

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Hub struct {
	Upgrader    *websocket.Upgrader
	Clients     map[*websocket.Conn]bool
	Register    chan *websocket.Conn
	Boardcaster chan *websocket.Conn
	Unregister  chan *websocket.Conn
}

func HubRun(h *Hub) {
	for {
		select {
		case client := <-h.Register:
			h.Clients[client] = true
			log.Println(h.Clients)
		case client := <-h.Unregister:
			delete(h.Clients, client)
			log.Println(h.Clients)
		}
	}
}

func NewHub() *Hub {
	return &Hub{
		Upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
		Register: make(chan *websocket.Conn),
		Clients:  make(map[*websocket.Conn]bool),
	}
}
