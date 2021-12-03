package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

type RequestMiddleware struct {
}

func (r *RequestMiddleware) InitRequest(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "X-Requested-With,X_Requested_With,content-type,token,action")
	param := make(map[string]interface{})
	raw, _ := c.GetRawData()
	if len(raw) == 0 {
		c.Request.ParseForm()
		for k, v := range c.Request.Form {
			param[k] = v[0]
		}
	} else {
		json.Unmarshal(raw, &param)
	}
	json.Unmarshal(raw, &param)
	if c.GetHeader("token") != "" {
		//TODO 临时
		param["login_uid"] = 1
	}
	c.Set("param", param)
}
