package controllers

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/web/ws"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type DrawingController struct {
	RoomService    core.RoomService
	DrawingService core.DrawingService
}

// Draw handles WebSocket connections for a specific room.
// After upgrading the HTTP connection to a WebSocket, it manages client registration.
func (dc DrawingController) Draw() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestedRoomID := ctx.Param("id")
		roomID := resolveRoomID(ctx)

		if _, err := dc.RoomService.CreateRoom(roomID); err != nil {
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

		if err := dc.RoomService.JoinRoom(roomID, participant); err != nil {
			log.Printf("failed to join room: %v", err)
			_ = participant.Close()
			return
		}

		defer func() {
			if err := dc.RoomService.LeaveRoom(roomID, participant); err != nil {
				log.Printf("failed to leave room: %v", err)
			}
			_ = participant.Close()
		}()

		participant.ReadMessage(ctx.Request.Context(), func(message []byte) {
			if err := dc.RoomService.BroadcastToRoom(roomID, string(message)); err != nil {
				log.Printf("failed to broadcast message: %v", err)
			}
			if err := dc.DrawingService.Enqueue(roomID, message); err != nil {
				log.Printf("failed to enqueue drawing message: %v", err)
			}
		})
	}
}

func createParticipant(ctx *gin.Context) (*ws.Participant, error) {
	hub := ws.NewWSHub()
	client, err := hub.NewClient(ctx.Writer, ctx.Request)
	if err != nil {
		return nil, err
	}
	participant := ws.NewParticipant(0, &client)

	return participant, nil
}

func resolveRoomID(ctx *gin.Context) string {
	pathRoomID := ctx.Param("id")
	if pathRoomID != "" {
		return pathRoomID
	}

	return fmt.Sprintf("%d", time.Now().UnixNano())
}

func (dc DrawingController) NotifyRoomCreated(participant core.Participant, roomID string) error {
	message, err := json.Marshal(ws.Message{
		Scope: string(ws.ScopeTypeLobby),
		Data: gin.H{
			"room_id": roomID,
		},
	})
	if err != nil {
		return err
	}

	participant.Notify(string(message))
	return nil
}
