package dao

import (
	"context"
	"github.com/feizhiwu/gs/kokomi"
	"github.com/jinzhu/gorm"
	"stage/app/model"
)

type UserDao struct {
	ctx      context.Context
	DB       *gorm.DB
	User     model.User
	UserList model.UserList
}

func User(ctx context.Context) *UserDao {
	return &UserDao{
		ctx: ctx,
		DB:  ctx.Value("main_db").(*gorm.DB),
	}
}

func (d *UserDao) Add() {
	d.DB.Create(&d.User)
}

func (d *UserDao) Update(data map[string]interface{}) {
	d.DB.Table("user").Where("id  = ?", data["id"]).Updates(data)
}

func (d *UserDao) GetOne() {
	d.DB.Where("id  = ?", d.User.Id).First(&d.User)
}

func (d *UserDao) Delete() {
	d.DB.Table("user").Delete(&d.User)
}

func (d *UserDao) GetAll(data map[string]interface{}) {
	db := d.DB.Model(model.User{})
	query := kokomi.NewQuery(&db, data)
	query.Like("name") //如果传参data["name"]，则进行like匹配查询
	query.List(&d.UserList.List).Pages(&d.UserList.Pages)
}
