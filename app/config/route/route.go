package route

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
