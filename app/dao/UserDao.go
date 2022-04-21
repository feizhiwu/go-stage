package dao

import (
	"context"
	"github.com/feizhiwu/gs/kokomi"
	"github.com/jinzhu/gorm"
	"stage/app/model"
)

type Base struct {
	ctx      context.Context
	mainDB   *gorm.DB
	clientDB *gorm.DB
}

type UserDao struct {
	Base
	User     model.User
	UserList model.UserList
}

func User(ctx context.Context) *UserDao {
	ud := new(UserDao)
	ud.ctx = ctx
	ud.mainDB = ctx.Value("main_db").(*gorm.DB)
	return ud
}

func (d *UserDao) Add() {
	d.mainDB.Create(&d.User)
}

func (d *UserDao) Update(data map[string]interface{}) {
	d.mainDB.Table("user").Where("id  = ?", data["id"]).Updates(data)
}

func (d *UserDao) GetOne() {
	d.mainDB.Where("id  = ?", d.User.Id).First(&d.User)
}

func (d *UserDao) Delete() {
	d.mainDB.Table("user").Delete(&d.User)
}

func (d *UserDao) GetAll(data map[string]interface{}) {
	db := d.mainDB.Model(model.User{})
	query := kokomi.Active(&db, data)
	query.Like("name") //如果传参data["name"]，则进行like匹配查询
	query.List(&d.UserList.List).Pages(&d.UserList.Pages)
}
