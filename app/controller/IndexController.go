package controller

import (
	"github.com/gin-gonic/gin"
	"stage/app/plugin/scene"
)

type IndexController struct {
	*scene.Drama
}

func Index(c *gin.Context) {
	s := &IndexController{
		scene.Action(c),
	}
	s.Show(10000)
}
