package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type RequestMiddleware struct {
}

func (r *RequestMiddleware) InitRequest(c *gin.Context) {
	param := make(map[string]interface{})
	raw, _ := c.GetRawData()
	json.Unmarshal(raw, &param)
	if c.GetHeader("token") != "" {
		//TODO 临时
		param["mid"] = 1
	}
	c.Set("param", param)
	c.Header("Access-Control-Allow-Origin", "*")
}
