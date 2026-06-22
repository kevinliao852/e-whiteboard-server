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

func (p *Participant) Notify(message string) {
	if p == nil || p.client == nil || *p.client == nil {
		return
	}

	if err := (*p.client).WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
		log.Println("WriteMessage error:", err)
	}
}

func (p *Participant) WriteMessage(data []byte) error {
	if p == nil || p.client == nil || *p.client == nil {
		return fmt.Errorf("websocket client is not initialized")
	}

	return (*p.client).WriteMessage(websocket.TextMessage, data)
}

func (p *Participant) Close() error {
	if p == nil || p.client == nil || *p.client == nil {
		return nil
	}

	return (*p.client).Close()
}

func (p *Participant) ReadMessage(
	ctx context.Context,
	onMessage func([]byte),
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
				return
			}

			fmt.Println("Received message from participant:", string(message))
			onMessage(message)
		}
	}
}
