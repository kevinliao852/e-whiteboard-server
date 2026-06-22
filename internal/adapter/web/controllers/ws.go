package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/web/authstate"
	"github.com/kevinliao852/e-whiteboard-server/internal/adapter/web/ws"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
)

type DrawingController struct {
	RoomService    core.RoomService
	DrawingService core.DrawingService
}

type PointResponse struct {
	ID           int    `json:"id"`
	WhiteboardID int    `json:"whiteboard_id"`
	Start        [2]int `json:"start"`
	End          [2]int `json:"end"`
}

type drawingIncomingMessage struct {
	Scope string          `json:"scope"`
	Data  json.RawMessage `json:"data"`
}

type cursorPayload struct {
	X      int  `json:"x"`
	Y      int  `json:"y"`
	Active bool `json:"active"`
}

type cursorEvent struct {
	Scope string          `json:"scope"`
	Data  cursorEventData `json:"data"`
}

type cursorEventData struct {
	ConnectionID string `json:"connection_id"`
	SenderID     uint   `json:"sender_id"`
	SenderName   string `json:"sender_name"`
	X            int    `json:"x,omitempty"`
	Y            int    `json:"y,omitempty"`
	Active       bool   `json:"active"`
}

// Draw handles WebSocket connections for a specific room.
// After upgrading the HTTP connection to a WebSocket, it manages client registration.
func (dc DrawingController) Draw() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomID := ctx.Param("id")
		if _, err := strconv.Atoi(roomID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid whiteboard id"})
			return
		}

		identity, ok := authstate.FromContext(ctx)
		if !ok {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
			return
		}

		displayName := identity.DisplayName
		if displayName == "" {
			displayName = "Unknown"
		}

		connectionID := fmt.Sprintf("%d", time.Now().UnixNano())

		if _, err := dc.RoomService.CreateRoom(roomID); err != nil {
			log.Printf("failed to create or load room: %v", err)
			return
		}

		participant, err := createParticipant(ctx, uint(identity.UserID))
		if err != nil {
			log.Printf("failed to create new participant: %v", err)
			return
		}

		if err := dc.RoomService.JoinRoom(roomID, participant); err != nil {
			log.Printf("failed to join room: %v", err)
			_ = participant.Close()
			return
		}

		defer func() {
			cursorLeave := cursorEvent{
				Scope: string(ws.ScopeTypeCursor),
				Data: cursorEventData{
					ConnectionID: connectionID,
					SenderID:     uint(identity.UserID),
					SenderName:   displayName,
					Active:       false,
				},
			}
			if payload, err := json.Marshal(cursorLeave); err == nil {
				if err := dc.RoomService.BroadcastToRoom(roomID, string(payload)); err != nil {
					log.Printf("failed to broadcast cursor leave: %v", err)
				}
			}
			if err := dc.RoomService.LeaveRoom(roomID, participant); err != nil {
				log.Printf("failed to leave room: %v", err)
			}
			_ = participant.Close()
		}()

		participant.ReadMessage(ctx.Request.Context(), func(message []byte) {
			var incoming drawingIncomingMessage
			if err := json.Unmarshal(message, &incoming); err != nil {
				log.Printf("failed to parse drawing message: %v", err)
				return
			}

			switch incoming.Scope {
			case string(ws.ScopeTypeWhiteboard):
				if err := dc.RoomService.BroadcastToRoom(roomID, string(message)); err != nil {
					log.Printf("failed to broadcast message: %v", err)
				}
				if err := dc.DrawingService.Enqueue(roomID, message); err != nil {
					log.Printf("failed to enqueue drawing message: %v", err)
				}
			case string(ws.ScopeTypeCursor):
				var payload cursorPayload
				if err := json.Unmarshal(incoming.Data, &payload); err != nil {
					log.Printf("failed to parse cursor payload: %v", err)
					return
				}

				cursor := cursorEvent{
					Scope: string(ws.ScopeTypeCursor),
					Data: cursorEventData{
						ConnectionID: connectionID,
						SenderID:     uint(identity.UserID),
						SenderName:   displayName,
						X:            payload.X,
						Y:            payload.Y,
						Active:       true,
					},
				}

				broadcast, err := json.Marshal(cursor)
				if err != nil {
					log.Printf("failed to marshal cursor event: %v", err)
					return
				}

				if err := dc.RoomService.BroadcastToRoom(roomID, string(broadcast)); err != nil {
					log.Printf("failed to broadcast cursor event: %v", err)
				}
			default:
				if err := dc.RoomService.BroadcastToRoom(roomID, string(message)); err != nil {
					log.Printf("failed to broadcast message: %v", err)
				}
			}
		})
	}
}

func (dc DrawingController) GetPoints(c *gin.Context) {
	whiteboardID, err := strconv.Atoi(c.Param("id"))
	if err != nil || whiteboardID <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid whiteboard id"})
		return
	}

	points, err := dc.DrawingService.ListCanvasData(whiteboardID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to load points"})
		return
	}

	response := make([]PointResponse, 0, len(points))
	for _, point := range points {
		response = append(response, PointResponse{
			ID:           point.ID,
			WhiteboardID: point.WhiteboardId,
			Start:        [2]int{point.StartX, point.StartY},
			End:          [2]int{point.EndX, point.EndY},
		})
	}

	c.JSON(http.StatusOK, response)
}

func createParticipant(ctx *gin.Context, userID uint) (*ws.Participant, error) {
	hub := ws.NewWSHub()
	client, err := hub.NewClient(ctx.Writer, ctx.Request)
	if err != nil {
		return nil, err
	}
	participant := ws.NewParticipant(userID, &client)

	return participant, nil
}
