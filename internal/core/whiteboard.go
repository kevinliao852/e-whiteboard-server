package core

import "time"

type Whiteboard struct {
	Id        uint
	UserId    uint
	Name      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

type WhiteboardModel interface {
	Create(wb *Whiteboard) error
	Delete(id uint) error
	GetByUserId(userId uint) ([]*Whiteboard, error)
}

type WhiteboardService interface {
	CreateWhiteboard(wb Whiteboard) (*Whiteboard, error)
	DeleteWhiteboard(id uint) error
	GetUserWhiteboards(userId uint) ([]*Whiteboard, error)
}
