package dao

import (
	"context"
	"github.com/feizhiwu/gs/kokomi"
	"stage/app/model"
)

type UserDao struct {
	BaseDao
	User     model.User
	UserList model.UserList
}

func User(ctx context.Context) *UserDao {
	dao := new(UserDao)
	dao.Active(ctx)
	return dao
}

func (d *UserDao) Add() {
	d.master.Create(&d.User)
}

func (d *UserDao) Update(data map[string]interface{}) {
	d.master.Table("user").Where("id  = ?", data["id"]).Updates(data)
}

func (d *UserDao) GetOne() {
	d.master.Where("id  = ?", d.User.Id).First(&d.User)
}

func (d *UserDao) Delete() {
	d.master.Table("user").Delete(&d.User)
}

func (d *UserDao) GetAll(data map[string]interface{}) {
	db := d.master.Model(model.User{})
	query := kokomi.Active(&db, data)
	query.Like("name") //如果传参data["name"]，则进行like匹配查询
	query.List(&d.UserList.List).Pages(&d.UserList.Pages)
}
