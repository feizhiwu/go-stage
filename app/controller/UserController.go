package controller

import (
	"github.com/feizhiwu/gs/albedo"
	"github.com/gin-gonic/gin"
	"stage/app/plugin/scene"
	"stage/app/service"
)

type UserController struct {
	*scene.Drama
	us *service.UserService
}

// User 控制器入口
func User(c *gin.Context) {
	s := &UserController{
		scene.Action(c),
		service.User(c.Request.Context()),
	}
	s.Get(s.list)
	s.Get(s.info)
	s.Post(s.add)
	s.Put(s.update)
	s.Delete(s.delete)
	s.Run()
}

func (s *UserController) add() {
	val := map[int]string{
		20001: "name",
		20002: "password",
	}
	s.Validate(val, s.Params)
	s.us.Add(s.Params)
	data := map[string]uint{
		"id": s.us.UD.User.Id,
	}
	s.Show(data)
}

func (s *UserController) list() {
	val := map[int]string{
		80007: "page",
	}
	s.Validate(val, s.Params)
	s.us.GetList(s.Params)
	s.Show(s.us.UD.UserList)
}

func (s *UserController) info() {
	s.HasKey(s.Params)
	s.us.GetInfo(albedo.MakeUint(s.Params["id"]))
	s.Show(s.us.UD.User)
}

func (s *UserController) update() {
	s.HasKey(s.Params)
	s.us.Update(s.Params)
	s.Show(scene.StatusOK)
}

func (s *UserController) delete() {
	s.HasKey(s.Params)
	s.us.Delete(albedo.MakeUint(s.Params["id"]))
	s.Show(scene.StatusOK)
}
