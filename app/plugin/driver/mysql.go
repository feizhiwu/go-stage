package driver

import (
	"context"
	"fmt"
	"github.com/feizhiwu/gs/albedo"
	"github.com/jinzhu/gorm"
	"stage/config/conf"
)

var (
	DBC map[string]*gorm.DB //db connect
	DBS []string            //dbs
)

// Connect 连接数据库
func Connect() {
	DBC = make(map[string]*gorm.DB)
	dbConf := conf.Config("db")
	for k, v := range dbConf.(map[interface{}]interface{}) {
		if _, ok := v.(map[interface{}]interface{}); !ok {
			connect("db", dbConf.(map[interface{}]interface{}))
			break
		} else {
			connect(albedo.MakeString(k), v.(map[interface{}]interface{}))
		}
	}
}

func connect(name string, options map[interface{}]interface{}) {
	ms := make(map[string]string)
	for k, v := range options {
		ms[albedo.MakeString(k)] = albedo.MakeString(v)
	}
	if ms["charset"] == "" {
		ms["charset"] = "utf8"
	}
	if ms["hostport"] == "" {
		ms["hostport"] = "3306"
	}
	db, err := gorm.Open(ms["datatype"],
		ms["username"]+":"+ms["password"]+"@tcp("+ms["hostname"]+":"+ms["hostport"]+")/"+ms["database"]+"?charset="+ms["charset"])
	if err != nil {
		panic(fmt.Sprintf("\x1b[31;20m[ERROR] %s\x1b[0m\n", err.Error()))
	}
	//全局禁用表复数
	db.SingularTable(true)
	//保存连接信息，键名与配置信息同名
	DBS = append(DBS, name)
	DBC[name] = db
}

func Begin(ctx context.Context) context.Context {
	var db *gorm.DB
	for _, v := range DBS {
		db = ctx.Value(v).(*gorm.DB).Begin()
		ctx = context.WithValue(ctx, v, db)
	}
	return ctx
}

func Commit(ctx context.Context) {
	for _, v := range DBS {
		ctx.Value(v).(*gorm.DB).Commit()
	}
}
