package service

import (
	"encoding/json"
	"toutGin/app/common"
	"toutGin/app/dao"
	"toutGin/app/model"
)

type UserService struct {
}

func (s *UserService) Add(data map[string]interface{}) uint {
	userDao := new(dao.UserDao)
	params := common.CopyParams([]string{"name", "password"}, data)
	json.Unmarshal(common.MakeJson(params), &userDao.User)
	return userDao.Add()
}

func (s *UserService) GetInfo(id uint) model.User {
	userDao := new(dao.UserDao)
	userDao.User.Id = id
	return userDao.GetOne()
}

func (s *UserService) Update(data map[string]interface{}) {
	userDao := new(dao.UserDao)
	params := common.CopyParams([]string{"id", "name", "password"}, data)
	userDao.Update(params)
}

func (s *UserService) Delete(id uint) {
	userDao := new(dao.UserDao)
	userDao.User.Id = id
	userDao.Delete()
}

func (s *UserService) GetList(data map[string]interface{}) interface{} {
	userDao := new(dao.UserDao)
	return userDao.GetAll(data)
}