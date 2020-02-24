package controller

import (
	"github.com/gin-gonic/gin"
	"toutGin/app/common"
	"toutGin/app/service"
)

type UserController struct {
	display *common.Display
	data    map[string]interface{}
}

//控制器入口
func (s *UserController) Run(c *gin.Context) {
	//防止变量污染
	s = new(UserController)
	s.display = &common.Display{Context: c}
	s.data = common.GetData(c)
	defer s.display.CatchPanic()
	switch {
	case c.Request.Method == "POST":
		s.display.IsLogin(s.data)
		s.add()
	case c.Request.Method == "GET":
		if s.data["id"] != nil {
			s.info()
		} else {
			s.list()
		}
	case c.Request.Method == "PUT":
		s.display.IsLogin(s.data)
		s.update()
	case c.Request.Method == "DELETE":
		s.display.IsLogin(s.data)
		s.delete()
	default:
		s.display.Show(common.StatusOK)
	}
	s.display.Finish()
}

func (s *UserController) add() {
	val := map[int]string{
		20001: "name",
		20002: "password",
	}
	s.display.IsEmpty(val, s.data)
	body := make(map[string]uint)
	body["id"] = new(service.UserService).Add(s.data)
	s.display.Show(body)
}

func (s *UserController) list() {
	val := map[int]string{
		80007: "page",
	}
	s.display.IsEmpty(val, s.data)
	body := new(service.UserService).GetList(s.data)
	s.display.Show(body)
}

func (s *UserController) info() {
	s.display.HasKey(s.data)
	body := new(service.UserService).GetInfo(common.MakeUint(s.data["id"]))
	s.display.Show(body)
}

func (s *UserController) update() {
	s.display.HasKey(s.data)
	new(service.UserService).Update(s.data)
	s.display.Show(common.StatusOK)
}

func (s *UserController) delete() {
	s.display.HasKey(s.data)
	new(service.UserService).Delete(common.MakeUint(s.data["id"]))
	s.display.Show(common.StatusOK)
}
