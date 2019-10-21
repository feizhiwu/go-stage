package service

import (
	"encoding/json"
	"toutGin/app/common"
	"toutGin/app/dao"
	"toutGin/app/model"
)

type UserService struct {
	UD *dao.UserDao
}

func (s *UserService) Add(data map[string]interface{}) uint {
	s.UD = new(dao.UserDao)
	params := common.CopyParams([]string{"name", "password"}, data)
	json.Unmarshal(common.MakeJson(params), &s.UD.User)
	return s.UD.Add()
}

func (s *UserService) GetInfo(id uint) model.User {
	s.UD = new(dao.UserDao)
	s.UD.User.Id = id
	return s.UD.User
}

func (s *UserService) Update(data map[string]interface{}) {
	s.UD = new(dao.UserDao)
	params := common.CopyParams([]string{"id", "name", "password"}, data)
	s.UD.Update(params)
}

func (s *UserService) Delete(id uint) {
	s.UD = new(dao.UserDao)
	s.UD.User.Id = id
	s.UD.Delete()
}

func (s *UserService) GetList(data map[string]interface{}) []model.User {
	s.UD = new(dao.UserDao)
	s.UD.GetAll(data)
	return s.UD.UserList
}
