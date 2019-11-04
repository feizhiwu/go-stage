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
	default:
		display.Show(common.StatusOk)
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
