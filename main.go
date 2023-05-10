package main

import (
	"github.com/gin-gonic/gin"
	"stage/app/plugin/driver"
	"stage/config/route"
)

func main() {
	engine := gin.New()
	driver.Connect()
	route.Run(engine)
	engine.Run(":8080")
}
