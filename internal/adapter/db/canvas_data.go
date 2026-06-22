package db

import (
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/kevinliao852/e-whiteboard-server/internal/database"
)

type WhiteboardCanvasData struct {
	ID           uint
	WhiteboardId uint `json:"whiteboard_id" binding:"required"`
	StartX       int  `json:"start_x" binding:"required"`
	EndX         int  `json:"end_x" binding:"required"`
	StartY       int  `json:"start_y" binding:"required"`
	EndY         int  `json:"end_y" binding:"required"`
	CreatedAt    time.Time
	UpdateAt     time.Time
}

func CreateCanvasData(data *core.CanvasData) error {
	row := &WhiteboardCanvasData{
		WhiteboardId: uint(data.WhiteboardId),
		StartX:       data.StartX,
		EndX:         data.EndX,
		StartY:       data.StartY,
		EndY:         data.EndY,
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
			StartX:       row.StartX,
			EndX:         row.EndX,
			StartY:       row.StartY,
			EndY:         row.EndY,
			CreatedAt:    row.CreatedAt,
			UpdateAt:     row.UpdateAt,
		})
	}

	return result, nil
}
