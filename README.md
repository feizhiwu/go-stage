# toutGin
个人风格的golang版RESTful API项目结构，gin+gorm，简单易上手

结构清晰简单，代码洁癖患者的福音
## 项目结构
```
|-app
    |-common 公共方法
        |-display.go 统一json格式输出
    |-controller 控制器
    |-dao 负责curd的
    |-middleware 中间件
    |-model 模型
    |-service 核心业务处理
|-config 配置文件和统一路由管理
    |-route
        |-route.go 路由配置文件
    |-message.yml 状态码配置文件
    |-config.go 配置方法
|-main.go 程序执行入口
```
## 模块调用流程
```
controller -> service -> dao
# controller严禁复杂业务，严禁直接调用dao，更严禁写sql语句
# 要优雅，不要凌乱，乖~
```
## REST URL
```
GET http://localhost:8080/v1/user
POST http://localhost:8080/v1/user
PUT http://localhost:8080/v1/user
DELETE http://localhost:8080/v1/user
# api 只接收 Content-Type: application/json 的传参，包括GET请求方式
```
## JSON RESULT
```
{
    "status": 10000,
    "msg": "请求成功",
    "body": null
}
```
## ROUTE 示例
```
package config

import (
	"github.com/gin-gonic/gin"
	"toutGin/app/controller"
	"toutGin/app/middleware"
)

type Route struct {
	Engine *gin.Engine
}

//路由路口
func (r *Route) Run() {
	r.Engine.Use(new(middleware.RequestMiddleware).InitRequest)
	r.index()
	r.v1()
}

func (r *Route) index() {
	r.Engine.Any("", new(controller.IndexController).Run)
}

func (r *Route) v1() {
	v1 := r.Engine.Group("v1")
	{
		v1.Any("/user", new(controller.UserController).Run)
	}
}
```
## CONTROLLER 示例
```
package controller

import (
	"github.com/gin-gonic/gin"
	"toutGin/app/common"
	"toutGin/app/service"
)

type UserController struct {
	data map[string]interface{}
}

//控制器入口
func (s *UserController) Run(c *gin.Context) {
	display = &common.Display{Context: c}
	data = common.GetData(c)
	defer display.CatchPanic()
	switch {
	case c.Request.Method == "POST":
		display.IsLogin(data)
		s.add()
	case c.Request.Method == "GET":
		if data["id"] != nil {
			s.info()
		} else {
			s.list()
		}
	case c.Request.Method == "PUT":
		display.IsLogin(data)
		s.update()
	case c.Request.Method == "DELETE":
		display.IsLogin(data)
		s.delete()
	}
}

func (s *UserController) add() {
	val := map[int]string{
		20001: "name",
		20002: "password",
	}
	display.IsEmpty(val, data)
	body := make(map[string]uint)
	body["id"] = new(service.UserService).Add(data)
	display.Show(body)
}

func (s *UserController) list() {
	val := map[int]string{
		80007: "page",
	}
	display.IsEmpty(val, data)
	body := new(service.UserService).GetList(data)
	display.Show(body)
}

func (s *UserController) info() {
	display.HasKey(data)
	body := new(service.UserService).GetInfo(common.MakeUint(data["id"]))
	display.Show(body)
}

func (s *UserController) update() {
	display.HasKey(data)
	new(service.UserService).Update(data)
	display.Show(common.StatusOk)
}

func (s *UserController) delete() {
	display.HasKey(data)
	new(service.UserService).Delete(common.MakeUint(data["id"]))
	display.Show(common.StatusOk)
}
```
## SERVICE 示例
```
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
	s.UD.Add()
	return s.UD.User.Id
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
```
## DAO 示例
```
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

func (d *UserDao) Add() {
	table := config.DB.Table("user")
	table.Create(&d.User)
	table.Last(&d.User)
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
```
## 接口统一返回
```
package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	StatusInit = 0
	StatusOk   = 10000
)

type Body struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Body   interface{} `json:"body"`
}

type Display struct {
	Context *gin.Context
	Body    Body
}

//统一输出api数据
func (d *Display) Show(mix interface{}) {
	message := new(Message)
	if d.Body.Status == StatusInit {
		d.Body.Status = StatusOk
	}
	if d.Body.Status != StatusOk {
		d.Body.Msg = message.GetMessage(d.Body.Status)
		d.Body.Body = nil
	} else {
		if val, ok := mix.(int); ok {
			d.Body.Status = val
			d.Body.Msg = message.GetMessage(d.Body.Status)
			d.Body.Body = nil
		} else if val, ok := mix.(string); ok {
			d.Body.Status = 11000
			d.Body.Msg = val
			d.Body.Body = nil
		} else {
			d.Body.Msg = message.GetMessage(d.Body.Status)
			d.Body.Body = mix
		}
	}
	d.Context.JSON(http.StatusOK, d.Body)
	d.Context.Abort()
}
```
## 测试user表结构
```
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```