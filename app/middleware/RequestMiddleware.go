package middleware

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"stage/app/plugin/scene"
)

func Request(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Credentials", "true")
	c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Header("Access-Control-Allow-Headers", "X-Requested-With,X_Requested_With,content-type,token,action")
	params := make(map[string]interface{})
	c.Request.ParseForm()
	raw, _ := c.GetRawData()
	if len(raw) == 0 {
		for k, v := range c.Request.Form {
			params[k] = v[0]
		}
	} else {
		json.Unmarshal(raw, &params)
	}
	json.Unmarshal(raw, &params)
	token := c.GetHeader("token")
	if token != "" || true {
		//coding 根据token获取用户信息
		params["login_uid"] = 1
	}
	c.Set("params", params)
	scene.Access(c, params)
	if !exclude(c) {
		if params["login_uid"] == nil {
			panic(80003)
		}
	}
}

//路由免登排除
func exclude(c *gin.Context) bool {
	rs := map[string][]string{
		"/": nil,
	}
	var action string
	for k, v := range rs {
		if k == c.Request.URL.Path {
			if v == nil {
				return true
			}
			action = c.GetHeader("action")
			for _, a := range v {
				if a == action {
					return true
				}
			}
		}
	}
	return false
}
