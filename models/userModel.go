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
