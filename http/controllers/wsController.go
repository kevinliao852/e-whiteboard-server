package controllers

import (
	"encoding/json"
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

		if _, ok := rooms[roomId]; !ok {
			room := wshub.NewRoom(roomId)
			rooms[roomId] = room
			go wshub.RunRoom(room)
			log.Println("Created new room: ", roomId)
		}

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
			mt, message, err := c.ReadMessage()
			if err != nil {
				log.Print("read:", err)
				break
			}
			//log.Printf("recv: %s", message)

			for client := range currentRoom.Clients {
				// broadcast data to everyone
				(*client).WriteMessage(mt, message)
				channel <- message
			}
		}
	}
}

type StoreWhiteboardMessage struct {
	RoomId string `json:"room_id"`
}

type Message struct {
	Start []uint `json:"start"`
	End   []uint `json:"end"`
}

func (swm *StoreWhiteboardMessage) SaveMessage(message []byte) error {
	// store message to db
	// example msg {"start":[305,245],"end":[312,245]}
	var parsedMsg Message
	err := json.Unmarshal(message, &parsedMsg)

	if err != nil {
		log.Println("Error parsing message", err, string(message))
		return err
	}

	whiteboardId, err := strconv.ParseUint(swm.RoomId, 10, 32)

	if err != nil {
		log.Println("Error parsing room id ", err, swm.RoomId)
		return err
	}

	uintWhiteboardId := uint(whiteboardId)

	data := &models.WhiteboardCanvasData{
		StartX:       parsedMsg.Start[0],
		StartY:       parsedMsg.Start[1],
		EndX:         parsedMsg.End[0],
		EndY:         parsedMsg.End[1],
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

	// use wshub.StoreMessageToDB to save message to db
	wshub.StoreMessageToDB(messageChannel, &storeWhiteboardMessage)
	defer close(messageChannel)
}
