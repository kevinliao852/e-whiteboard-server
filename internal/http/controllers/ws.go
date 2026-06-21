package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/kevinliao852/e-whiteboard-server/internal/wsmodel"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type DrawingController struct {
	RoomService core.RoomService
}

// Draw handles WebSocket connections for a specific room.
// After upgrading the HTTP connection to a WebSocket, it manages client registration,
func (dc DrawingController) Draw() gin.HandlerFunc {
	var rooms = sync.Map{}
	roomCtx := context.Background()

	return func(ctx *gin.Context) {
		requestedRoomID := ctx.Param("id")
		roomID := resolveRoomID(ctx)

		currentRoom, err := dc.createRoom(roomCtx, &rooms, roomID)
		if err != nil {
			log.Printf("failed to create or load room: %v", err)
			return
		}

		participant, err := createParticipant(ctx)
		if err != nil {
			log.Printf("failed to create new participant: %v", err)
			return
		}

		if requestedRoomID == "" {
			if err := dc.NotifyRoomCreated(participant, roomID); err != nil {
				log.Printf("failed to notify created room: %v", err)
				_ = participant.Close()
				return
			}
		}

		currentRoom.Register <- participant
		defer func() {
			currentRoom.Unregister <- participant
			_ = participant.Close()
		}()

		participant.ReadMessage(ctx, currentRoom.Broadcast)
	}
}

func createParticipant(ctx *gin.Context) (*wsmodel.Participant, error) {

	hub := wsmodel.NewWSHub()
	client, err := hub.NewClient(ctx.Writer, ctx.Request)
	if err != nil {
		return nil, err
	}
	participant := wsmodel.NewParticipant(0, &client)

	return participant, nil
}

func resolveRoomID(ctx *gin.Context) string {
	pathRoomID := ctx.Param("id")
	if pathRoomID != "" {
		return pathRoomID
	}

	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (dc DrawingController) createRoom(ctx context.Context, rooms *sync.Map, id string) (*wsmodel.Room, error) {
	newRoom := wsmodel.NewRoom(id)

	actual, loaded := rooms.LoadOrStore(id, newRoom)

	r, ok := actual.(*wsmodel.Room)
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

func (dc DrawingController) NotifyRoomCreated(participant *wsmodel.Participant, roomID string) error {
	message, err := json.Marshal(wsmodel.Message{
		Scope: string(wsmodel.ScopeTypeLobby),
		Data: gin.H{
			"room_id": roomID,
		},
	})
	if err != nil {
		return err
	}

	return participant.WriteMessage(message)
}
