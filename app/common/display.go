package common

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stage/config/conf"
)

const (
	StatusInit = 0
	StatusOK   = 10000
	StatusFail = 11000
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

	jsonApi struct {
		Context *gin.Context
		Body    Body
	}
)

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

func NewDisplay(c *gin.Context) *Display {
	d := &Display{
		Context: c,
		Params:  GetParams(c),
		funcs:   make(map[string]func()),
	}
	return d
}

type Display struct {
	*gin.Context
	Render
	Params map[string]interface{}
	status int
	funcs  map[string]func()
	checks struct {
		verify  bool
		actions []string
	}
}

func (d *Display) Get(f func()) {
	d.funcs["GET-"+GetFuncName(f)] = f
}

func (d *Display) Post(f func()) {
	d.funcs["POST-"+GetFuncName(f)] = f
}

func (d *Display) Put(f func()) {
	d.funcs["PUT-"+GetFuncName(f)] = f
}

func (d *Display) Delete(f func()) {
	d.funcs["DELETE-"+GetFuncName(f)] = f
}

// Show 统一输出api数据
func (d *Display) Show(mix interface{}) {
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
func (d *Display) Validate(val map[int]string, data map[string]interface{}) {
	d.status = StatusOK
	for k, v := range val {
		if data[v] == nil {
			panic(k)
		}
	}
}

// HasKey 检测更新主键是否为空
func (d *Display) HasKey(data map[string]interface{}) {
	if data["id"] == nil {
		panic(80001)
	}
}

// ForceLogin 需要强制登陆
func (d *Display) ForceLogin() {
	if GetParams(d.Context)["login_uid"] == nil {
		panic(80003)
	}
}

// CheckAction 行为检查
func (d *Display) CheckAction(value string) bool {
	d.status = StatusWarn
	action := d.Context.GetHeader("action")
	if action != "" && action == value {
		d.status = StatusOK
		return true
	}
	return false
}

func (d *Display) Run() {
	action := d.GetHeader("action")
	f := d.funcs[d.Request.Method+"-"+action]
	if f != nil {
		if len(d.checks.actions) > 0 {
			if d.checks.verify && InArray(len(d.checks.actions), func(i int) bool {
				return d.checks.actions[i] == action
			}) {
				d.ForceLogin()
				f()
			} else if !d.checks.verify && !InArray(len(d.checks.actions), func(i int) bool {
				return d.checks.actions[i] == action
			}) {
				d.ForceLogin()
				f()
			} else {
				f()
			}
		} else {
			f()
		}
	} else {
		d.Show(StatusWarn)
	}
}

// CatchPanic 统一中断输出
func (d *Display) CatchPanic() {
	if r := recover(); r != nil {
		d.Show(r)
	}
}
