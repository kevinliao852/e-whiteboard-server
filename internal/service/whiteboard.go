package service

import "github.com/kevinliao852/e-whiteboard-server/internal/core"

type WhiteboardSVC struct {
	CreateFn      func(wb *core.Whiteboard) error
	DeleteFn      func(id uint) error
	GetByUserIDFn func(userId uint) ([]*core.Whiteboard, error)
}

func NewWhiteboardSVC(createFn func(wb *core.Whiteboard) error, deleteFn func(id uint) error, getByUserIDFn func(userId uint) ([]*core.Whiteboard, error)) *WhiteboardSVC {
	return &WhiteboardSVC{
		CreateFn:      createFn,
		DeleteFn:      deleteFn,
		GetByUserIDFn: getByUserIDFn,
	}
}

func (s *WhiteboardSVC) CreateWhiteboard(wb core.Whiteboard) (*core.Whiteboard, error) {
	if err := s.CreateFn(&wb); err != nil {
		return nil, err
	}

	return &wb, nil
}

func (s *WhiteboardSVC) GetUserWhiteboards(userId uint) ([]*core.Whiteboard, error) {
	return s.GetByUserIDFn(userId)
}

func (s *WhiteboardSVC) DeleteWhiteboard(whiteboardId uint) error {
	return s.DeleteFn(whiteboardId)
}

var _ core.WhiteboardService = (*WhiteboardSVC)(nil)
