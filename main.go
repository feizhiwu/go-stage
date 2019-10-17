package main

import (
	"github.com/gin-gonic/gin"
	"toutGin/app/common"
	"toutGin/app/config"
)

func main() {
	gin.SetMode(gin.DebugMode)
	new(common.Database).GetConnect()
	defer common.DB.Close()
	engine := gin.New()
	(&config.Route{Engine: engine}).Run()
	engine.Run(":8080")
}
