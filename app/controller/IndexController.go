package controller

import (
	"github.com/gin-gonic/gin"
	"toutGin/app/common"
)

type IndexController struct {
}

func (s *IndexController) Run(c *gin.Context) {
	display = &common.Display{Context: c}
	display.Show("来啦，老弟")
}
