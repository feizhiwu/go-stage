package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
)

func Init(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "X-Requested-With,X_Requested_With,content-type,token,action")
	params := make(map[string]interface{})
	raw, _ := c.GetRawData()
	if len(raw) == 0 {
		c.Request.ParseForm()
		for k, v := range c.Request.Form {
			params[k] = v[0]
		}
	} else {
		json.Unmarshal(raw, &params)
	}
	json.Unmarshal(raw, &params)
	if c.GetHeader("token") != "" {
		//TODO 临时
		params["login_uid"] = 1
	}
	c.Set("params", params)
}
