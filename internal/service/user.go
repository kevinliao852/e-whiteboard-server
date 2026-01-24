package service

import "github.com/kevinliao852/e-whiteboard-server/internal/core"

type UserSVC struct {
	Model core.UserModel
}

func (svc *UserSVC) GetUser(id string) (*core.User, error) {
	return svc.Model.GetById(id)
}

func (svc *UserSVC) Register(user *core.User) error {
	return svc.Model.Create(user)
}

func (svc *UserSVC) GetUserByGoogleId(gid string) (*core.User, error) {
	return svc.Model.GetByGoogleId(gid)
}
