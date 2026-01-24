package core

import "time"

type User struct {
	ID          int
	DisplayName string
	Email       string
	GoogleID    string
	CreateAt    time.Time
	UpdateAt    time.Time
}

type UserModel interface {
	Create(user *User) error
	GetById(id string) (*User, error)
	GetByGoogleId(gid string) (*User, error)
}

type UserService interface {
	GetUser(id string) (*User, error)
	Register(user *User) error
	GetUserByGoogleId(gid string) (*User, error)
}
