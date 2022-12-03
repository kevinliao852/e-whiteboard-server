package models

import (
	"app/database"
	"time"
)

type Whiteboard struct {
	ID        uint
	UserId    uint `json:"user_id" binding:"required"`
	CreatedAt time.Time
	UpdateAt  time.Time
}

func CreateAWhiteboard(wb *Whiteboard) error {
	if err := database.DB.Create(wb).Error; err != nil {
		return err
	}

	return nil
}

func DeleteAWhiteboard(id uint) error {

	if err := database.DB.Where("id = ?", id).Delete(&Whiteboard{}).Error; err != nil {
		return err
	}
	return nil
}

func GetWhiteboardsById(wbs *[]Whiteboard, id uint) error {
	if err := database.DB.Find(wbs, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func GetWhiteboardsByUserId(wbs *[]Whiteboard, userId uint) error {
	if err := database.DB.Find(wbs, "user_id = ?", userId).Error; err != nil {
		return err
	}
	return nil
}