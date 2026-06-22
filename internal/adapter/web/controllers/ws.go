package controllers

import (
	"net/http"
	"strconv"

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

// Draw handles WebSocket connections for a specific room.
// After upgrading the HTTP connection to a WebSocket, it manages client registration.
func (dc DrawingController) Draw() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		roomID := ctx.Param("id")
		if _, err := strconv.Atoi(roomID); err != nil {
			ctx.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "invalid whiteboard id"})
			return
		}

		if _, err := dc.RoomService.CreateRoom(roomID); err != nil {
			log.Printf("failed to create or load room: %v", err)
			return
		}

		participant, err := createParticipant(ctx)
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

func createParticipant(ctx *gin.Context) (*ws.Participant, error) {
	hub := ws.NewWSHub()
	client, err := hub.NewClient(ctx.Writer, ctx.Request)
	if err != nil {
		return nil, err
	}
	participant := ws.NewParticipant(0, &client)

	return participant, nil
}
