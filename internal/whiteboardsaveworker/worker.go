package whiteboardsaveworker

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

	"github.com/kevinliao852/e-whiteboard-server/internal/model"
	"github.com/kevinliao852/e-whiteboard-server/internal/wsmodel"
)

type WhiteboardSaveWorker struct {
	RoomId string `json:"room_id"`
}

type WhiteboardMessage struct {
	Start [2]uint `json:"start"`
	End   [2]uint `json:"end"`
}

// SaveMessage saves a whiteboard message to the database.
func (swm *WhiteboardSaveWorker) SaveMessage(message []byte) error {
	var parsedMsg wsmodel.Message
	if err := json.Unmarshal(message, &parsedMsg); err != nil {
		log.Println("[SaveMessage] Error parsing message", err, string(message))
		return err
	}

	if parsedMsg.Scope != string(wsmodel.ScopeTypeWhiteboard) {
		return fmt.Errorf("[SaveMessage] unsupported scope: %s", parsedMsg.Scope)
	}

	whiteboardId, err := strconv.ParseUint(swm.RoomId, 10, 32)
	if err != nil {
		return fmt.Errorf("[SaveMessage] Error parsing room id:%v, %+v", swm.RoomId, whiteboardId)
	}

	data, err := json.Marshal(parsedMsg.Data)
	if err != nil {
		return fmt.Errorf("[SaveMessage] failed to marshal message payload: %w", err)
	}

	var wmd WhiteboardMessage
	if err = json.Unmarshal(data, &wmd); err != nil {
		return fmt.Errorf("[SaveMessage] WhiteboardMessage Unmarshal failed %+v", err)
	}

	if err = model.Create(&model.WhiteboardCanvasData{
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

	var storeWhiteboardMessage WhiteboardSaveWorker
	storeWhiteboardMessage.RoomId = roomId

	errChan := make(chan error)

	// use wshub.StoreMessageToDB to save message to db
	wsmodel.StartMessagePersistenceWorker(
		context.Background(),
		messageChannel, &storeWhiteboardMessage, &errChan)

	go func() {
		for c := range errChan {
			log.Println("Error saving message ", c)
		}
	}()

	return messageChannel
}

func ParseMessage(rawMessage []byte) (wsmodel.Message, error) {

	var message wsmodel.Message
	parseErr := json.Unmarshal(rawMessage, &message)

	if parseErr != nil {
		log.Println("Error parsing message", parseErr, string(rawMessage))
	}
	return message, nil
}
