package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	"github.com/kevinliao852/e-whiteboard-server/internal/models"
	"github.com/kevinliao852/e-whiteboard-server/internal/wshub"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

// WebsocketRoute handles WebSocket connections for a specific room.
// After upgrading the HTTP connection to a WebSocket, it manages client registration,
func WebsocketRoute() gin.HandlerFunc {
	var rooms = sync.Map{}
	roomCtx := context.Background()

	return func(ctx *gin.Context) {
		id := ctx.Param("id")
		currentRoom, err := createRoom(roomCtx, &rooms, id)
		if err != nil {
			log.Printf("failed to create or load room: %v", err)
			return
		}

		participant, err := createParticipant(ctx, currentRoom)
		if err != nil {
			log.Printf("failed to create new participant: %v", err)
			return
		}

		currentRoom.Register <- participant
		defer func() {
			currentRoom.Unregister <- participant
			_ = participant.Close()
		}()

		participant.ReadMessage(ctx, currentRoom.Broadcast)
	}
}

func createParticipant(
	ctx *gin.Context,
	room *wshub.Room,
) (*wshub.Participant, error) {

	hub := wshub.NewWSHub()
	client, err := hub.NewClient(ctx.Writer, ctx.Request)
	if err != nil {
		return nil, err
	}
	participant := wshub.NewParticipant(0, &client)

	return participant, nil
}

func createRoom(ctx context.Context, rooms *sync.Map, id string) (*wshub.Room, error) {
	newRoom := wshub.NewRoom(id)
	actual, loaded := rooms.LoadOrStore(id, newRoom)

	r, ok := actual.(*wshub.Room)
	if !ok {
		return nil, fmt.Errorf("invalid room type for id=%s", id)
	}

	if !loaded {
		go func() {
			r.Run(ctx)
			rooms.Delete(id)
		}()
		log.Printf("Created and started new room with id=%s", id)
	}

	return r, nil
}

func ParseMessage(rawMessage []byte) (wshub.Message, error) {

	var message wshub.Message
	parseErr := json.Unmarshal(rawMessage, &message)

	if parseErr != nil {
		log.Println("Error parsing message", parseErr, string(rawMessage))
	}
	return message, nil
}

type StoreWhiteboardMessage struct {
	RoomId string `json:"room_id"`
}

type WhiteboardMessage struct {
	Start [2]uint `json:"start"`
	End   [2]uint `json:"end"`
}

// SaveMessage saves a whiteboard message to the database.
func (swm *StoreWhiteboardMessage) SaveMessage(message []byte) error {
	var parsedMsg wshub.Message
	if err := json.Unmarshal(message, &parsedMsg); err != nil {
		log.Println("[SaveMessage] Error parsing message", err, string(message))
		return err
	}

	whiteboardId, err := strconv.ParseUint(swm.RoomId, 10, 32)
	if err != nil {
		return fmt.Errorf("[SaveMessage] Error parsing room id:%v, %+v", swm.RoomId, whiteboardId)
	}

	var wmd WhiteboardMessage
	if err = json.Unmarshal(message, &wmd); err != nil {
		return fmt.Errorf("[SaveMessage] WhiteboardMessage Unmarshal failed %+v", err)
	}

	if err = models.Create(&models.WhiteboardCanvasData{
		StartX:       wmd.Start[0],
		StartY:       wmd.Start[1],
		EndX:         wmd.End[0],
		EndY:         wmd.End[1],
		WhiteboardId: uint(whiteboardId),
	}); err != nil {
		log.Println("Error saving message ", err)
		return err
	}

	return nil
}

func NewWhiteboardSaveWorker(roomId string) chan []byte {
	messageChannel := make(chan []byte, 100)

	var storeWhiteboardMessage StoreWhiteboardMessage
	storeWhiteboardMessage.RoomId = roomId

	errChan := make(chan error)

	// use wshub.StoreMessageToDB to save message to db
	wshub.StartMessagePersistenceWorker(
		context.Background(),
		messageChannel, &storeWhiteboardMessage, &errChan)

	go func() {
		for c := range errChan {
			log.Println("Error saving message ", c)
		}
	}()

	return messageChannel
}
