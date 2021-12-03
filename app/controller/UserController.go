package controller

import (
	"github.com/gin-gonic/gin"
	"stage/app/common"
	"stage/app/service"
)

type UserController struct {
	display *common.Display
	data    map[string]interface{}
	us      *service.UserService
}

// User 控制器入口
func User(c *gin.Context) {
	s := &UserController{
		display: &common.Display{Context: c},
		data:    common.GetData(c),
		us:      new(service.UserService),
	}
	defer s.display.CatchPanic()
	switch {
	case c.Request.Method == "GET":
		if s.display.CheckAction("list") {
			s.list()
		} else if s.display.CheckAction("info") {
			s.info()
		}
	case c.Request.Method == "POST":
		if s.display.CheckAction("add") {
			s.add()
		}
	case c.Request.Method == "PUT":
		if s.display.CheckAction("update") {
			s.update()
		}
	case c.Request.Method == "DELETE":
		if s.display.CheckAction("delete") {
			s.delete()
		}
	}
	s.display.Finish()
}

func (s *UserController) add() {
	val := map[int]string{
		20001: "name",
		20002: "password",
	}
	s.display.IsEmpty(val, s.data)
	s.us.Add(s.data)
	data := map[string]uint{
		"id": s.us.UD.User.Id,
	}
	s.display.Show(data)
}

func (s *UserController) list() {
	val := map[int]string{
		80007: "page",
	}
	s.display.IsEmpty(val, s.data)
	s.us.GetList(s.data)
	s.display.Show(s.us.UD.UserList)
}

func (s *UserController) info() {
	s.display.HasKey(s.data)
	s.us.GetInfo(common.MakeUint(s.data["id"]))
	s.display.Show(s.us.UD.User)
}

func (s *UserController) update() {
	s.display.HasKey(s.data)
	s.us.Update(s.data)
	s.display.Show(common.StatusOK)
}

func (s *UserController) delete() {
	s.display.HasKey(s.data)
	s.us.Delete(common.MakeUint(s.data["id"]))
	s.display.Show(common.StatusOK)
}
