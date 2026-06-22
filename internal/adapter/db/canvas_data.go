package db

import (
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
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

func (w *WhiteboardCanvasData) Create(data *core.CanvasData) error {
	w.WhiteboardId = uint(data.WhiteboardId)
	w.StartX = uint(data.StartX)
	w.EndX = uint(data.EndX)
	w.StartY = uint(data.StartY)
	w.EndY = uint(data.EndY)

	if err := database.DB.Create(w).Error; err != nil {
		return err
	}

	return nil
}

func Create(wb *WhiteboardCanvasData) error {
	if err := database.DB.Create(wb).Error; err != nil {
		return err
	}

	return nil
}

var _ core.CanvasDataInterface = (*WhiteboardCanvasData)(nil)
