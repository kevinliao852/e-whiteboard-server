package core

import "time"

type CanvasData struct {
	ID           int
	WhiteboardId int
	StartX       int
	EndX         int
	StartY       int
	EndY         int
	CreatedAt    time.Time
	UpdateAt     time.Time
}

type CanvasDataInterface interface {
	Create(data *CanvasData) error
	GetById(id int) (*CanvasData, error)
}
