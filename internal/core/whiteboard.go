package core

import "time"

type Whiteboard struct {
	Id        uint
	UserId    uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WhiteboardInterface interface {
	Create(wb *Whiteboard) error
	Delete(id uint) error
	GetById(id uint) ([]Whiteboard, error)
	GetByUserId(userId uint) ([]Whiteboard, error)
	Update(wb *Whiteboard) error
}
