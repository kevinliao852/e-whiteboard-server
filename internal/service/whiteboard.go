package service

import "github.com/kevinliao852/e-whiteboard-server/internal/core"

type WhiteboardService struct {
	model core.WhiteboardModel
}

func NewCreateWhiteboardService() *WhiteboardService {
	return &WhiteboardService{}
}

func (s *WhiteboardService) CreateWhiteboard(wb core.Whiteboard) (*core.Whiteboard, error) {
	s.model.Create(&wb)
	return &wb, nil
}

func (s *WhiteboardService) GetUserWhiteboards(userId uint) ([]*core.Whiteboard, error) {
	return s.model.GetByUserId(userId)
}

func (s *WhiteboardService) DeleteWhiteboard(whiteboardId uint) error {
	return s.model.Delete(whiteboardId)
}
