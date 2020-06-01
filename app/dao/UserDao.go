package dao

import (
	"toutGin/app/common"
	"toutGin/app/model"
	"toutGin/config"
)

type UserDao struct {
	User     model.User
	UserList []model.User
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
	limit := 20
	config.DB.Table("user").Limit(limit).Offset(common.GetOffset(data["page"], limit)).Find(&d.UserList)
}
