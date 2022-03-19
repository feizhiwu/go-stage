package main

import (
	"github.com/gin-gonic/gin"
	"stage/config/conf"
	"stage/config/route"
)

func main() {
	engine := gin.New()
	conf.ConnectDB()
	route.Run(engine)
	engine.Run(":8080")
}
