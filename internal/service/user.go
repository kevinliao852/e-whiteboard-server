package service

import "github.com/kevinliao852/e-whiteboard-server/internal/core"

type UserSVC struct {
	CreateFn        func(user *core.User) error
	GetByIDFn       func(id string) (*core.User, error)
	GetByGoogleIDFn func(gid string) (*core.User, error)
}

func (svc *UserSVC) GetUser(id string) (*core.User, error) {
	return svc.GetByIDFn(id)
}

func (svc *UserSVC) Register(user *core.User) error {
	return svc.CreateFn(user)
}

func (svc *UserSVC) GetUserByGoogleId(gid string) (*core.User, error) {
	return svc.GetByGoogleIDFn(gid)
}
