package models

import (
	"app/database"
	"time"
)

type Whiteboard struct {
	Id        uint
	UserId    uint `json:"user_id" binding:"required"`
	UserId    uint `json:"user-id" binding:"required"`
	Name      string
	CreatedAt time.Time
	UpdateAt  time.Time
}

func (w *Whiteboard) CreateAWhiteboard(wb *Whiteboard) error {
	if err := database.DB.Create(wb).Error; err != nil {
		return err
	}

	return nil
}

func (w *Whiteboard) DeleteAWhiteboard(id uint) error {

	if err := database.DB.Where("id = ?", id).Delete(&Whiteboard{}).Error; err != nil {
		return err
	}
	return nil
}

func (w *Whiteboard) GetWhiteboardsById(wbs *[]Whiteboard, id uint) error {
	if err := database.DB.Find(wbs, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func (w *Whiteboard) GetWhiteboardsByUserId(wbs *[]Whiteboard, userId uint) error {
	if err := database.DB.Find(wbs, "user_id = ?", userId).Error; err != nil {
		return err
	}
	return nil
}

func (w *Whiteboard) UpdateAWhiteboard(wb *Whiteboard) error {
	if err := database.DB.Save(wb).Error; err != nil {
		return err
	}
	return nil
}
