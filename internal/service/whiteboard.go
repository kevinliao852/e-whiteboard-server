package service

import "github.com/kevinliao852/e-whiteboard-server/internal/core"

type WhiteboardSVC struct {
	Model core.WhiteboardModel
}

func NewCreateWhiteboardService() *WhiteboardSVC {
	return &WhiteboardSVC{}
}

func (s *WhiteboardSVC) CreateWhiteboard(wb core.Whiteboard) (*core.Whiteboard, error) {
	if err := s.Model.Create(&wb); err != nil {
		return nil, err
	}

	return &wb, nil
}

func (s *WhiteboardSVC) GetUserWhiteboards(userId uint) ([]*core.Whiteboard, error) {
	return s.Model.GetByUserId(userId)
}

func (s *WhiteboardSVC) DeleteWhiteboard(whiteboardId uint) error {
	return s.Model.Delete(whiteboardId)
}

var _ core.WhiteboardService = (*WhiteboardSVC)(nil)
