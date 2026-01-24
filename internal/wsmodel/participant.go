package wsmodel

import (
	"context"
	"fmt"
	"log"

	"github.com/gorilla/websocket"
)

type Client = *websocket.Conn

type Participant struct {
	UserID uint
	client *Client
}

func NewParticipant(userID uint, client *Client) *Participant {
	return &Participant{
		UserID: userID,
		client: client,
	}
}

func (p *Participant) WriteMessage(data []byte) error {
	return (*p.client).WriteMessage(websocket.TextMessage, data)
}

func (p *Participant) Close() error {
	return (*p.client).Close()
}

func (p *Participant) ReadMessage(
	ctx context.Context,
	ch chan []byte,
) {
	for {
		select {
		case <-ctx.Done():
			log.Println("Context done, stopping ReadMessage for participant")
			return
		default:
			_, message, err := (*p.client).ReadMessage()
			if err != nil {
				log.Println("ReadMessage error:", err)
				log.Println("Closing participant connection")
				break
			}

			fmt.Println("Received message from participant:", string(message))
			ch <- message
		}
	}
}
