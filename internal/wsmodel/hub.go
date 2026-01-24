package wsmodel

import (
	"net/http"

	"github.com/gorilla/websocket"
)

type WSHub struct {
	upgrader *websocket.Upgrader
}

func NewWSHub() *WSHub {
	return &WSHub{
		upgrader: &websocket.Upgrader{
			CheckOrigin: func(r *http.Request) bool {
				return true

			},
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
		},
	}
}

type WSHuber interface {
	NewClient(
		writer http.ResponseWriter,
		request *http.Request,
	) (Client, error)
}

var _ WSHuber = (*WSHub)(nil)

func (wsh *WSHub) NewClient(writer http.ResponseWriter, request *http.Request) (Client, error) {
	return wsh.upgrader.Upgrade(writer, request, nil)
}
