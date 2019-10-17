package dao

import (
	"toutGin/app/common"
	"toutGin/app/model"
)

type UserDao struct {
	User model.User
}

func (d *UserDao) Add() uint {
	table := common.DB.Table("user")
	table.Create(&d.User)
	table.Last(&d.User)
	return d.User.Id
}

func (d *UserDao) Update(data map[string]interface{}) {
	common.DB.Table("user").Where("id  = ?", data["id"]).Updates(data)
}

func (d *UserDao) GetOne() model.User {
	common.DB.Table("user").Where("id  = ?", d.User.Id).First(&d.User)
	return d.User
}

func (d *UserDao) Delete() {
	common.DB.Table("user").Delete(&d.User)
}

func (d *UserDao) GetAll(data map[string]interface{}) []model.User {
	var users []model.User
	limit := 20
	common.DB.Table("user").Limit(limit).Offset(common.GetOffset(data["page"], limit)).Find(&users)
	return users
}
