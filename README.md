# toutGin
个人风格的golang版RESTful API项目结构，gin+gorm，简单易上手

结构清晰简单，代码洁癖患者的福音
## 获取源码
git clone https://github.com/feizhiwu/toutGin.git
## 项目结构
```
|-app
    |-common 公共方法
        |-display.go 统一json格式输出
        |-function.go 公共函数
    |-controller 控制器
    |-dao 负责curd的
    |-middleware 中间件
    |-model 模型
    |-service 核心业务处理
|-config 配置文件和统一路由管理
    |-route
        |-route.go 路由配置文件
    |-message.yml 状态码配置文件
    |-config.go 配置方法
|-main.go 程序执行入口
```
## 模块调用流程
```
controller -> service -> dao
# controller严禁复杂业务，严禁直接调用dao，更严禁写sql语句
```
## RESTful API
```
GET http://localhost:8080/v1/user
POST http://localhost:8080/v1/user
PUT http://localhost:8080/v1/user
DELETE http://localhost:8080/v1/user
# api POST，PUT，DELETE 推荐使用 body json 传参，GET兼容 body 和 url 传参
```
## JSON RESULT
```
{
    "status": 10000,
    "msg": "请求成功",
    "body": null
}
```
## ROUTE 示例
```
type Route struct {
	Engine *gin.Engine
}

//路由路口
func (r *Route) Run() {
	r.Engine.Use(new(middleware.RequestMiddleware).InitRequest)
	r.index()
	r.v1()
}

func (r *Route) index() {
	r.Engine.Any("", new(controller.IndexController).Run)
}

func (r *Route) v1() {
	v1 := r.Engine.Group("v1")
	{
		v1.Any("/user", controller.User)
	}
}
```
## CONTROLLER 示例
```
type UserController struct {
	display *common.Display
	data    map[string]interface{}
	us      *service.UserService
}

//控制器入口
func User(c *gin.Context) {
	s := &UserController{
		display: &common.Display{Context: c},
		data:    common.GetData(c),
		us:      new(service.UserService),
	}
	defer s.display.CatchPanic()
	switch {
	case c.Request.Method == "POST":
		s.display.IsLogin(s.data)
		s.add()
	case c.Request.Method == "GET":
		if s.data["id"] != nil {
			s.info()
		} else {
			s.list()
		}
	case c.Request.Method == "PUT":
		s.display.IsLogin(s.data)
		s.update()
	case c.Request.Method == "DELETE":
		s.display.IsLogin(s.data)
		s.delete()
	s.display.Finish()
}

func (s *UserController) add() {
	val := map[int]string{
		20001: "name",
		20002: "password",
	}
	s.display.IsEmpty(val, s.data)
	s.us.Add(s.data)
	data := map[string]uint{
		"id": s.us.UD.User.Id,
	}
	s.display.Show(data)
}

func (s *UserController) list() {
	val := map[int]string{
		80007: "page",
	}
	s.display.IsEmpty(val, s.data)
	s.us.GetList(s.data)
	s.display.Show(s.us.UD.UserList)
}

func (s *UserController) info() {
	s.display.HasKey(s.data)
	s.us.GetInfo(common.MakeUint(s.data["id"]))
	s.display.Show(s.us.UD.User)
}

func (s *UserController) update() {
	s.display.HasKey(s.data)
	s.us.Update(s.data)
	s.display.Show(common.StatusOK)
}

func (s *UserController) delete() {
	s.display.HasKey(s.data)
	s.us.Delete(common.MakeUint(s.data["id"]))
	s.display.Show(common.StatusOK)
}
```
## SERVICE 示例
```
type UserService struct {
	UD *dao.UserDao
}

func (s *UserService) Add(data map[string]interface{}) {
	s.UD = new(dao.UserDao)
	params := common.CopyParams([]string{"name", "password"}, data)
	json.Unmarshal(common.MakeJson(params), &s.UD.User)
	s.UD.User.Password = common.EncryptPass(s.UD.User.Password)
	s.UD.Add()
}

func (s *UserService) GetInfo(id uint) {
	s.UD = new(dao.UserDao)
	s.UD.User.Id = id
	s.UD.GetOne()
}

func (s *UserService) Update(data map[string]interface{}) {
	s.UD = new(dao.UserDao)
	params := common.CopyParams([]string{"id", "name", "password"}, data)
	s.UD.Update(params)
}

func (s *UserService) Delete(id uint) {
	s.UD = new(dao.UserDao)
	s.UD.User.Id = id
	s.UD.Delete()
}

func (s *UserService) GetList(data map[string]interface{}) {
	s.UD = new(dao.UserDao)
	s.UD.GetAll(data)
}
```
## DAO 示例
```
type UserDao struct {
	User     model.User
	UserList []model.User
}

func (d *UserDao) Add() {
	config.DB.Create(&d.User)
}

func (d *UserDao) Update(data map[string]interface{}) {
	config.DB.Table("user").Where("id  = ?", data["id"]).Updates(data)
}

func (d *UserDao) GetOne() {
	config.DB.Where("id  = ?", d.User.Id).First(&d.User)
}

func (d *UserDao) Delete() {
	config.DB.Table("user").Delete(&d.User)
}

func (d *UserDao) GetAll(data map[string]interface{}) {
	limit := 20
	config.DB.Table("user").Limit(limit).Offset(common.GetOffset(data["page"], limit)).Find(&d.UserList)
}
```
## 接口统一返回
```
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

func (d *Display) CheckAction(data map[string]interface{}, value string) bool {
	d.Status = StatusWarn
	if data["action"].(string) == value {
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
```
## 测试user表结构
```
CREATE TABLE `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(255) DEFAULT NULL,
  `password` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
```