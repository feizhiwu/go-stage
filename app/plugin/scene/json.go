package scene

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stage/config/conf"
)

type jsonApi struct {
	Context *gin.Context
	Body    Body
}

func Json(c *gin.Context) *jsonApi {
	return &jsonApi{Context: c}
}

func (j *jsonApi) Output(mix interface{}) {
	j.Body.Status = StatusOK
	if val, ok := mix.(int); ok {
		j.Body.Status = val
		j.Body.Msg = conf.Message(j.Body.Status)
		j.Body.Body = nil
	} else if val, ok := mix.(string); ok {
		j.Body.Status = StatusFail
		j.Body.Msg = val
		j.Body.Body = nil
	} else {
		j.Body.Msg = conf.Message(j.Body.Status)
		j.Body.Body = mix
	}
	if j.Body.Msg == "" {
		j.Body.Msg = conf.Message(StatusFail)
	}
	j.Context.JSON(http.StatusOK, j.Body)
	j.Context.Abort()
}
