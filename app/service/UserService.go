package service

import (
	"context"
	"encoding/json"
	"github.com/feizhiwu/gs/albedo"
	"stage/app/dao"
	"stage/app/plugin"
	"stage/app/plugin/driver"
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
	//启动事务
	s.ctx = driver.Begin(s.ctx)
	s.UD = dao.User(s.ctx)
	params := plugin.CopyParams([]string{"name", "password"}, data)
	json.Unmarshal(albedo.MakeJson(params), &s.UD.User)
	s.UD.User.Password = plugin.EncryptPass(s.UD.User.Password)
	s.UD.Add()
	//提交事务
	driver.Commit(s.ctx)
}

func (s *UserService) GetInfo(id uint) {
	s.UD = dao.User(s.ctx)
	s.UD.User.Id = id
	s.UD.GetOne()
}

func (s *UserService) Update(data map[string]interface{}) {
	s.UD = dao.User(s.ctx)
	params := plugin.CopyParams([]string{"id", "name", "password"}, data)
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
