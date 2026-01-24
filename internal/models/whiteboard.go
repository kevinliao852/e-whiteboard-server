package models

import (
	"time"

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

// CreateAWhiteboard creates a new whiteboard record in the database.
func (w *Whiteboard) CreateAWhiteboard(wb *Whiteboard) error {
	if err := database.DB.Create(wb).Error; err != nil {
		return err
	}

	return nil
}

// DeleteAWhiteboard deletes a whiteboard record from the database by its ID.
func (w *Whiteboard) DeleteAWhiteboard(id uint) error {
	if err := database.DB.Where("id = ?", id).Delete(&Whiteboard{}).Error; err != nil {
		return err
	}
	return nil
}

// GetWhiteboardsById retrieves whiteboard records from the database by their ID.
func (w *Whiteboard) GetWhiteboardsById(wbs *[]Whiteboard, id uint) error {
	if err := database.DB.Find(wbs, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

// GetWhiteboardsByUserId retrieves all whiteboard records for a specific user.
func (w *Whiteboard) GetWhiteboardsByUserId(wbs *[]Whiteboard, userId uint) error {
	if err := database.DB.Find(wbs, "user_id = ?", userId).Error; err != nil {
		return err
	}
	return nil
}

// UpdateAWhiteboard updates an existing whiteboard record in the database.
func (w *Whiteboard) UpdateAWhiteboard(wb *Whiteboard) error {
	if err := database.DB.Save(wb).Error; err != nil {
		return err
	}
	return nil
}
