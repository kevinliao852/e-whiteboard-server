package models

import (
	"app/database"
	"time"
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

func CreateAWhiteboardCanvasData(wb *WhiteboardCanvasData) error {
	if err := database.DB.Create(wb).Error; err != nil {
		return err
	}

	return nil
}

func DeleteAWhiteboardCanvasData(id uint) error {

	if err := database.DB.Where("id = ?", id).Delete(&WhiteboardCanvasData{}).Error; err != nil {
		return err
	}
	return nil
}

func GetWhiteboardCanvasDataById(wb *WhiteboardCanvasData, id uint) error {
	if err := database.DB.Find(wb, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}
