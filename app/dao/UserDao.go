package dao

import (
	"toutGin/app/common"
	"toutGin/app/config"
	"toutGin/app/model"
)

type UserDao struct {
	User     model.User
	UserList []model.User
}

func (d *UserDao) Add() uint {
	table := config.DB.Table("user")
	table.Create(&d.User)
	table.Last(&d.User)
	return d.User.Id
}

func (d *UserDao) Update(data map[string]interface{}) {
	config.DB.Table("user").Where("id  = ?", data["id"]).Updates(data)
}

func (d *UserDao) GetOne() {
	config.DB.Table("user").Where("id  = ?", d.User.Id).First(&d.User)
}

func (d *UserDao) Delete() {
	config.DB.Table("user").Delete(&d.User)
}

func (d *UserDao) GetAll(data map[string]interface{}) {
	limit := 20
	config.DB.Table("user").Limit(limit).Offset(common.GetOffset(data["page"], limit)).Find(&d.UserList)
}
