package controller

import (
	"github.com/gin-gonic/gin"
	"stage/app/common"
)

type IndexController struct {
	*common.Display
}

func Index(c *gin.Context) {
	s := &IndexController{
		common.NewDisplay(c),
	}
	s.Show(10000)
}
