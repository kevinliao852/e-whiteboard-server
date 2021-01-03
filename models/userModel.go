package models

import (
	"app/database"
)

type User struct {
	ID uint `json:"id"`
	Name string `json:"name"`
	Password string `json:"password"`
	Email string `json:"email"`
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