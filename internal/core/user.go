package core

import "time"

type User struct {
	Id          int
	DisplayName string
	Email       string
	GoogleId    string
	CreateAt    time.Time
	UpdateAt    time.Time
}

type UserInterface interface {
	Create(user *User) error
	GetById(id int) (*User, error)
	GetByGoogleId(gid string) (*User, error)
}
