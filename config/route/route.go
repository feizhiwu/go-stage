package route

import (
	"github.com/gin-gonic/gin"
	"stage/app/controller"
	"stage/app/middleware"
)

type route struct {
	Engine *gin.Engine
}

// Run 路由路口
func Run(engine *gin.Engine) {
	r := route{engine}
	r.Engine.Use(middleware.Init, middleware.Recovery)
	r.Engine.Any("", controller.Index)
	r.v1()
}
func (r *route) v1() {
	v1 := r.Engine.Group("v1")
	{
		v1.Any("/user", controller.User)
	}
}
