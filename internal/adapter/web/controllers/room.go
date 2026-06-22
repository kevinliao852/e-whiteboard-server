package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kevinliao852/e-whiteboard-server/internal/core"
)

type RoomController struct {
	service core.RoomService
}

type RoomResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Participants int    `json:"participants"`
	Activity     string `json:"activity"`
}

func NewRoomController(svc core.RoomService) *RoomController {
	return &RoomController{
		service: svc,
	}
}

func (ctrl *RoomController) ListRooms(c *gin.Context) {
	rooms := ctrl.service.ListRooms()
	response := make([]RoomResponse, 0, len(rooms))

	for _, room := range rooms {
		response = append(response, RoomResponse{
			ID:           room.ID,
			Name:         room.Name,
			Status:       room.Status,
			Participants: len(room.Participants),
			Activity:     formatActivity(room.LastActivity),
		})
	}

	c.JSON(http.StatusOK, response)
}

func formatActivity(lastActivity time.Time) string {
	if lastActivity.IsZero() {
		return "unknown"
	}

	elapsed := time.Since(lastActivity)
	switch {
	case elapsed < time.Minute:
		return "just now"
	case elapsed < time.Hour:
		return pluralDuration(int(elapsed.Minutes()), "m")
	case elapsed < 24*time.Hour:
		return pluralDuration(int(elapsed.Hours()), "h")
	default:
		return pluralDuration(int(elapsed.Hours()/24), "d")
	}
}

func pluralDuration(value int, unit string) string {
	return fmt.Sprintf("%d%s ago", value, unit)
}
