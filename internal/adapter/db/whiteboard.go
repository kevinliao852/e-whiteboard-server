package db

import (
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/kevinliao852/e-whiteboard-server/internal/database"
)

// Whiteboard represents the whiteboard model in the database.
type Whiteboard struct {
	Id        uint
	UserId    uint `json:"user-id" binding:"required"`
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time // corrected 'UpdateAt' to 'UpdatedAt'
}

func CreateWhiteboard(wb *core.Whiteboard) error {
	if err := database.DB.Create(wb).Error; err != nil {
		return err
	}

	return nil
}

func DeleteWhiteboard(id uint) error {
	if err := database.DB.Where("id = ?", id).Delete(&Whiteboard{}).Error; err != nil {
		return err
	}
	return nil
}

func GetWhiteboardsByUserID(userId uint) ([]*core.Whiteboard, error) {
	wbs := make([]Whiteboard, 0)
	if err := database.DB.Find(&wbs, "user_id = ?", userId).Error; err != nil {
		return nil, err
	}

	var result []*core.Whiteboard
	for i := range wbs {
		result = append(result, &core.Whiteboard{
			Id:        wbs[i].Id,
			UserId:    wbs[i].UserId,
			Name:      wbs[i].Name,
			CreatedAt: wbs[i].CreatedAt,
			UpdatedAt: wbs[i].UpdatedAt,
		})
	}

	return result, nil
}
