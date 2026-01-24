package models

import (
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/database"
)

type WhiteboardCanvasData struct {
	ID           uint
	WhiteboardId uint `json:"whiteboard_id" binding:"required"`
	StartX       uint `json:"start_x" binding:"required"`
	EndX         uint `json:"end_x" binding:"required"`
	StartY       uint `json:"start_y" binding:"required"`
	EndY         uint `json:"end_y" binding:"required"`
	CreatedAt    time.Time
	UpdateAt     time.Time
}

func Create(wb *WhiteboardCanvasData) error {
	if err := database.DB.Create(wb).Error; err != nil {
		return err
	}

	return nil
}
