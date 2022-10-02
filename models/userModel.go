package models

import (
	"app/database"
	"time"
)

type User struct {
	ID          uint
	DisplayName string `json:"display_name"`
	Email       string `json:"email" binding:"required" gorm:"unique"`
	GoogleId    string `json:"google_id" binding:"required"`
	CreateAt    time.Time
	UpdateAt    time.Time
}

func CreateAUser(user *User) error {
	if err := database.DB.Create(user).Error; err != nil {
		return err
	}
	return nil
}

func GetAllUsers(user *[]User) error {
	if err := database.DB.Find(user).Error; err != nil {

		return err
	}
	return nil
}

func DeleteAUser(name string) error {
	var user []User
	if getAllUsersErr := GetAllUsers(&user); getAllUsersErr != nil {
		return getAllUsersErr
	}

	if err := database.DB.Where("name = ?", name).Delete(user).Error; err != nil {
		return err
	}
	return nil
}

func GetUserById(user *User, id string) error {
	if err := database.DB.First(user, "id = ?", id).Error; err != nil {
		return err
	}
	return nil
}

func GetUserByGoogleId(user *User, gid string) error {
	if err := database.DB.Find(user, "google_id = ?", gid).Error; err != nil {
		return err
	}
	return nil
}
