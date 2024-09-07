package controllers

import (
	"encoding/json"
	"fmt"
	"strconv"

	"app/models"
	"app/wshub"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

func WebsocketRoute() gin.HandlerFunc {
	var rooms = make(map[string]*wshub.Room)
	return func(ctx *gin.Context) {

		roomId := ctx.Param("id")
		errChan := make(chan error)

		if _, ok := rooms[roomId]; !ok {
			room := wshub.NewRoom(roomId)
			rooms[roomId] = room
			go wshub.RunRoom(room, &errChan)
			log.Println("Created new room: ", roomId)
		}

		go func() {
			for c := range errChan {
				log.Println("Error saving message ", c)
			}
		}()

		currentRoom := rooms[roomId]
		c, err := currentRoom.Upgrader.Upgrade(ctx.Writer, ctx.Request, nil)

		if err != nil {
			log.Print("upgrade error:", err)
			return
		}

		channel := make(chan []byte)

		go WhiteboardSaveWorker(roomId, channel)
		log.Println("Created new WhiteboardSaveWorker client", c.RemoteAddr().String())
		defer c.Close()
		defer delete(currentRoom.Clients, c)
		currentRoom.Register <- c

		for {
			mt, rawMessage, err := c.ReadMessage()

			if err != nil {
				log.Print("read:", err)
				break
			}

			var message wshub.Message
			parseErr := json.Unmarshal(rawMessage, &message)

			switch message.Scope {
			case string(wshub.ScopeTypeWhiteboard):
				// TODO getWhiteboaryHistoryById
			}

			if parseErr != nil {
				log.Println("Error parsing message", parseErr, string(rawMessage))
				break
			}

			fmt.Println("Message scope: ", message.Scope, " data: ", message.Data)
			fmt.Println(string(rawMessage), parseErr)

			if message.Scope == string(wshub.ScopeTypeWhiteboard) {
				for client := range currentRoom.Clients {
					log.Println(currentRoom.Clients)

					if client == c {
						continue
					}

					err = client.WriteMessage(mt, rawMessage)

					if err != nil {
						log.Println("Error writing message", err)
					}

					err = (*client).WriteMessage(mt, rawMessage)

					if err != nil {
						log.Println("Error writing message", err)
					}

					channel <- rawMessage
				}
			}

		}
	}
}

func ParseMessage(rawMessage []byte) (interface{}, error) {

	var message wshub.Message
	parseErr := json.Unmarshal(rawMessage, &message)

	if parseErr != nil {
		log.Println("Error parsing message", parseErr, string(rawMessage))
	}
	return nil, nil
}

type StoreWhiteboardMessage struct {
	RoomId string `json:"room_id"`
}

type WhiteboardMessage struct {
	Start []uint `json:"start"`
	End   []uint `json:"end"`
}

func (swm *StoreWhiteboardMessage) SaveMessage(message []byte) error {
	// store message to db
	var parsedMsg wshub.Message
	err := json.Unmarshal(message, &parsedMsg)

	if err != nil {
		log.Println("[SaveMessage] Error parsing message", err, string(message))
		return err
	}

	whiteboardId, err := strconv.ParseUint(swm.RoomId, 10, 32)

	if err != nil {
		return fmt.Errorf("[SaveMessage] Error parsing room id:%v, %+v", swm.RoomId, whiteboardId)
	}

	uintWhiteboardId := uint(whiteboardId)

	fmt.Println(parsedMsg)

	marshalWmd, err := json.Marshal(parsedMsg.Data)

	if err != nil {
		return fmt.Errorf("[SaveMessage] WhiteboardMessage marchal failed %+v", err)
	}

	var wmd WhiteboardMessage

	err = json.Unmarshal(marshalWmd, &wmd)

	if err != nil {
		return fmt.Errorf("[SaveMessage] WhiteboardMessage Unmarshal failed %+v", err)
	}

	data := &models.WhiteboardCanvasData{
		StartX:       wmd.Start[0],
		StartY:       wmd.Start[1],
		EndX:         wmd.End[0],
		EndY:         wmd.End[1],
		WhiteboardId: uintWhiteboardId,
	}

	err = models.CreateAWhiteboardCanvasData(data)

	if err != nil {
		log.Println("Error saving message ", err)
		return err
	}

	return nil
}

func WhiteboardSaveWorker(roomId string, messageChannel chan []byte) {
	// save message to db

	// store whiteboard message
	var storeWhiteboardMessage StoreWhiteboardMessage
	storeWhiteboardMessage.RoomId = roomId

	errChan := make(chan error)

	// use wshub.StoreMessageToDB to save message to db
	wshub.StoreMessageToDB(messageChannel, &storeWhiteboardMessage, &errChan)
	defer close(messageChannel)

	go func() {
		for c := range errChan {
			log.Println("Error saving message ", c)
		}
	}()
}
