package main

import (
	"github.com/gin-gonic/gin"
	"stage/config/conf"
	"stage/config/route"
)

func main() {
	gin.SetMode(gin.DebugMode)
	new(conf.Database).GetConnect()
	defer conf.DB.Close()
	engine := gin.New()
	(&route.Route{Engine: engine}).Run()
	engine.Run(":8080")
}
