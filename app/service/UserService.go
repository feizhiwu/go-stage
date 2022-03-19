package service

import (
	"context"
	"encoding/json"
	"stage/app/common"
	"stage/app/dao"
)

type UserService struct {
	ctx context.Context
	UD  *dao.UserDao
}

func User(ctx context.Context) *UserService {
	return &UserService{
		ctx: ctx,
	}
}

func (s *UserService) Add(data map[string]interface{}) {
	s.UD = dao.User(s.ctx)
	params := common.CopyParams([]string{"name", "password"}, data)
	json.Unmarshal(common.MakeJson(params), &s.UD.User)
	s.UD.User.Password = common.EncryptPass(s.UD.User.Password)
	s.UD.Add()
}

func (s *UserService) GetInfo(id uint) {
	s.UD = dao.User(s.ctx)
	s.UD.User.Id = id
	s.UD.GetOne()
}

func (s *UserService) Update(data map[string]interface{}) {
	s.UD = dao.User(s.ctx)
	params := common.CopyParams([]string{"id", "name", "password"}, data)
	s.UD.Update(params)
}

func (s *UserService) Delete(id uint) {
	s.UD = dao.User(s.ctx)
	s.UD.User.Id = id
	s.UD.Delete()
}

func (s *UserService) GetList(data map[string]interface{}) {
	s.UD = dao.User(s.ctx)
	s.UD.GetAll(data)
}
