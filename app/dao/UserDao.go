package dao

import (
	"github.com/feizhiwu/gs/kokomi"
	"stage/app/model"
	"stage/config"
)

type UserDao struct {
	User     model.User
	UserList model.UserList
}

func (d *UserDao) Add() {
	config.DB.Create(&d.User)
}

func (d *UserDao) Update(data map[string]interface{}) {
	config.DB.Table("user").Where("id  = ?", data["id"]).Updates(data)
}

func (d *UserDao) GetOne() {
	config.DB.Where("id  = ?", d.User.Id).First(&d.User)
}

func (d *UserDao) Delete() {
	config.DB.Table("user").Delete(&d.User)
}

func (d *UserDao) GetAll(data map[string]interface{}) {
	db := config.DB.Model(model.User{})
	query := kokomi.NewQuery(&db, data)
	query.Like("name") //如果传参data["name"]，则进行like匹配查询
	query.List(&d.UserList.List).Pages(&d.UserList.Pages)
}
