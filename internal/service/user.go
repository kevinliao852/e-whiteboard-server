package service

import "github.com/kevinliao852/e-whiteboard-server/internal/core"

type UserService struct {
	model core.UserModel
}

func (svc *UserService) GetUser(id string) (*core.User, error) {
	return svc.model.GetById(id)
}

func (svc *UserService) Register(user *core.User) error {
	return svc.model.Create(user)
}

func (svc *UserService) GetUserByGoogleId(gid string) (*core.User, error) {
	return svc.model.GetByGoogleId(gid)
}
