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

func CreateCanvasData(data *core.CanvasData) error {
	row := &WhiteboardCanvasData{
		WhiteboardId: uint(data.WhiteboardId),
		StartX:       uint(data.StartX),
		EndX:         uint(data.EndX),
		StartY:       uint(data.StartY),
		EndY:         uint(data.EndY),
	}

	if err := database.DB.Create(row).Error; err != nil {
		return err
	}

	return nil
}

func GetCanvasDataByWhiteboardID(whiteboardID int) ([]core.CanvasData, error) {
	var rows []WhiteboardCanvasData
	if err := database.DB.
		Where("whiteboard_id = ?", whiteboardID).
		Order("id ASC").
		Find(&rows).Error; err != nil {
		return nil, err
	}

	result := make([]core.CanvasData, 0, len(rows))
	for _, row := range rows {
		result = append(result, core.CanvasData{
			ID:           int(row.ID),
			WhiteboardId: int(row.WhiteboardId),
			StartX:       int(row.StartX),
			EndX:         int(row.EndX),
			StartY:       int(row.StartY),
			EndY:         int(row.EndY),
			CreatedAt:    row.CreatedAt,
			UpdateAt:     row.UpdateAt,
		})
	}

	return result, nil
}
