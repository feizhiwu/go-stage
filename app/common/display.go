package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	StatusInit = 0
	StatusOk   = 10000
)

type Body struct {
	Status int         `json:"status"`
	Msg    string      `json:"msg"`
	Body   interface{} `json:"body"`
}

type Display struct {
	Context *gin.Context
	Body    Body
}

//统一输出api数据
func (d *Display) Show(mix interface{}) {
	message := new(Message)
	if d.Body.Status == StatusInit {
		d.Body.Status = StatusOk
	}
	if d.Body.Status != StatusOk {
		d.Body.Msg = message.GetMessage(d.Body.Status)
		d.Body.Body = nil
	} else {
		if val, ok := mix.(int); ok {
			d.Body.Status = val
			d.Body.Msg = message.GetMessage(d.Body.Status)
			d.Body.Body = nil
		} else if val, ok := mix.(string); ok {
			d.Body.Status = 11000
			d.Body.Msg = val
			d.Body.Body = nil
		} else {
			d.Body.Msg = message.GetMessage(d.Body.Status)
			d.Body.Body = mix
		}
	}
	d.Context.JSON(http.StatusOK, d.Body)
	d.Context.Abort()
}

//参数检测
func (d *Display) IsEmpty(val map[int]string, data map[string]interface{}) {
	d.Body.Status = StatusOk
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
	if data["mid"] == nil {
		panic(80003)
	}
}

//统一中断输出
func (d *Display) CatchPanic() {
	if r := recover(); r != nil {
		d.Show(r)
	}
}
