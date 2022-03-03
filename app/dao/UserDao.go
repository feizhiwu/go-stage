package dao

import (
	"github.com/feizhiwu/gs/kokomi"
	"stage/app/model"
	"stage/config/conf"
)

type UserDao struct {
	User     model.User
	UserList model.UserList
}

func (d *UserDao) Add() {
	conf.DB.Create(&d.User)
}

func (d *UserDao) Update(data map[string]interface{}) {
	conf.DB.Table("user").Where("id  = ?", data["id"]).Updates(data)
}

func (d *UserDao) GetOne() {
	conf.DB.Where("id  = ?", d.User.Id).First(&d.User)
}

func (d *UserDao) Delete() {
	conf.DB.Table("user").Delete(&d.User)
}

func (d *UserDao) GetAll(data map[string]interface{}) {
	db := conf.DB.Model(model.User{})
	query := kokomi.NewQuery(&db, data)
	query.Like("name") //如果传参data["name"]，则进行like匹配查询
	query.List(&d.UserList.List).Pages(&d.UserList.Pages)
}
