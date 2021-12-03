package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stage/config"
)

const (
	StatusInit = 0
	StatusOK   = 10000
	StatusWarn = 80006
)

type (
	Render interface {
		Output(interface{})
	}

	Body struct {
		Status int         `json:"status"`
		Msg    string      `json:"msg"`
		Body   interface{} `json:"body"`
	}

	JsonApi struct {
		Context *gin.Context
		Body    Body
	}
)

func (j *JsonApi) Output(mix interface{}) {
	message := new(config.Message)
	j.Body.Status = StatusOK
	if val, ok := mix.(int); ok {
		j.Body.Status = val
		j.Body.Msg = message.GetMessage(j.Body.Status)
		j.Body.Body = nil
	} else if val, ok := mix.(string); ok {
		j.Body.Status = 11000
		j.Body.Msg = val
		j.Body.Body = nil
	} else {
		j.Body.Msg = message.GetMessage(j.Body.Status)
		j.Body.Body = mix
	}
	j.Context.JSON(http.StatusOK, j.Body)
	j.Context.Abort()
}

type Display struct {
	Context *gin.Context
	Status  int
	Render  Render
}

//统一输出api数据
func (d *Display) Show(mix interface{}) {
	if d.Status == StatusInit {
		d.Status = StatusOK
	}
	//默认json格式
	if d.Render == nil {
		d.Render = &JsonApi{Context: d.Context}
	}
	d.Render.Output(mix)
}

//参数检测
func (d *Display) IsEmpty(val map[int]string, data map[string]interface{}) {
	d.Status = StatusOK
	for k, v := range val {
		if data[v] == nil {
			panic(k)
		}
	}
}

//检测更新主键是否为空
func (d *Display) HasKey(data map[string]interface{}) {
	if data["id"] == nil {
		panic(80001)
	}
}

//检测是否登录
func (d *Display) IsLogin(data map[string]interface{}) {
	if data["login_uid"] == nil {
		panic(80003)
	}
}

// CheckAction 行为检查
func (d *Display) CheckAction(value string) bool {
	d.Status = StatusWarn
	action := d.Context.GetHeader("action")
	if action != "" && action == value {
		d.Status = StatusOK
		return true
	}
	return false
}

func (d *Display) Finish() {
	if d.Status != StatusOK {
		d.Show(StatusWarn)
	}
}

//统一中断输出
func (d *Display) CatchPanic() {
	if r := recover(); r != nil {
		d.Show(r)
	}
}
