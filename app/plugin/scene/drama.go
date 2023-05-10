package scene

import (
	"github.com/gin-gonic/gin"
	"stage/app/plugin"
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
)

// Action 请开始你的表演
func Action(c *gin.Context) *Drama {
	d := &Drama{
		Context: c,
		Params:  plugin.GetParams(c),
		funcs:   make(map[string]func()),
	}
	return d
}

type Drama struct {
	*gin.Context
	Render
	Params map[string]interface{}
	status int
	funcs  map[string]func()
}

func (d *Drama) Get(f func()) {
	d.funcs["GET-"+plugin.GetFuncName(f)] = f
}

func (d *Drama) Post(f func()) {
	d.funcs["POST-"+plugin.GetFuncName(f)] = f
}

func (d *Drama) Put(f func()) {
	d.funcs["PUT-"+plugin.GetFuncName(f)] = f
}

func (d *Drama) Delete(f func()) {
	d.funcs["DELETE-"+plugin.GetFuncName(f)] = f
}

// Show 统一输出api数据
func (d *Drama) Show(mix interface{}) {
	if d.status == StatusInit {
		d.status = StatusOK
	}
	//默认json格式
	if d.Render == nil {
		d.Render = Json(d.Context)
	}
	d.Render.Output(mix)
}

// Validate 参数检测
func (d *Drama) Validate(val map[int]string, data map[string]interface{}) {
	d.status = StatusOK
	for k, v := range val {
		if data[v] == nil {
			panic(k)
		}
	}
}

// HasKey 检测更新主键是否为空
func (d *Drama) HasKey(data map[string]interface{}) {
	if data["id"] == nil {
		panic(80001)
	}
}

// ForceLogin 需要强制登陆
func (d *Drama) ForceLogin() {
	if plugin.GetParams(d.Context)["login_uid"] == nil {
		panic(80003)
	}
}

// CheckAction 行为检查
func (d *Drama) CheckAction(value string) bool {
	d.status = StatusWarn
	action := d.Context.GetHeader("action")
	if action != "" && action == value {
		d.status = StatusOK
		return true
	}
	return false
}

func (d *Drama) Run() {
	action := d.GetHeader("action")
	f := d.funcs[d.Request.Method+"-"+action]
	if f != nil {
		f()
	} else {
		d.Show(StatusWarn)
	}
}

// CatchPanic 统一中断输出
func (d *Drama) CatchPanic() {
	if r := recover(); r != nil {
		d.Show(r)
	}
}
