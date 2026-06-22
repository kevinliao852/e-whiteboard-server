package db

import (
	"time"

	"github.com/kevinliao852/e-whiteboard-server/internal/core"
	"github.com/kevinliao852/e-whiteboard-server/internal/database"
)

type User struct {
	Id          uint
	DisplayName string `json:"display_name"`
	Email       string `json:"email" binding:"required" gorm:"unique"`
	GoogleId    string `json:"google_id" binding:"required"`
	CreateAt    time.Time
	UpdateAt    time.Time
}

func CreateUser(user *core.User) error {
	if err := database.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByGoogleId(gid string) (*core.User, error) {
	var user User
	if err := database.DB.First(&user, "google_id = ?", gid).Error; err != nil {
		return nil, err
	}
	return &core.User{
		ID:          int(user.Id),
		DisplayName: user.DisplayName,
		Email:       user.Email,
		GoogleID:    user.GoogleId,
		CreateAt:    user.CreateAt,
		UpdateAt:    user.UpdateAt,
	}, nil
}

func GetUserByID(id string) (*core.User, error) {
	var user User
	if err := database.DB.First(&user, "id = ?", id).Error; err != nil {
		return nil, err
	}

	return &core.User{
		ID:          int(user.Id),
		DisplayName: user.DisplayName,
		Email:       user.Email,
		GoogleID:    user.GoogleId,
		CreateAt:    user.CreateAt,
		UpdateAt:    user.UpdateAt,
	}, nil
}
