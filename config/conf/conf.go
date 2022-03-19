package conf

import (
	"fmt"
	"github.com/feizhiwu/gs/albedo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	DBC map[string]*gorm.DB
	DBS []string
)

func Config(key string) interface{} {
	dir, _ := os.Getwd()
	filePath := path.Join(dir, "/config/config.yml")
	fileData, _ := ioutil.ReadFile(filePath)
	var config map[interface{}]interface{}
	yaml.Unmarshal(fileData, &config)
	config = config[gin.Mode()].(map[interface{}]interface{})
	keys := strings.Split(key, ".")
	length := len(keys)
	if length == 1 {
		return config[key]
	} else {
		var value interface{}
		for _, v := range keys {
			if value == nil {
				value = config[v]
			} else {
				value = value.(map[interface{}]interface{})[v]
			}
		}
		return value
	}
}

func Message(status int) string {
	var msg map[int]string
	var filePath string
	dir, _ := os.Getwd()
	filePath = path.Join(dir, "/config/message.yml")
	fileData, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(fileData, &msg)
	return msg[status]
}

func ConnectDB() {
	DBC = make(map[string]*gorm.DB)
	dbConf := Config("db")
	if len(dbConf.(map[interface{}]interface{})) > 1 {
		for k, v := range dbConf.(map[interface{}]interface{}) {
			connectDB(albedo.MakeString(k), v.(map[interface{}]interface{}))
		}
	} else {
		connectDB("db", dbConf.(map[interface{}]interface{}))
	}
}

func connectDB(name string, options map[interface{}]interface{}) {
	ms := make(map[string]string)
	for k, v := range options {
		ms[albedo.MakeString(k)] = albedo.MakeString(v)
	}
	if ms["charset"] == "" {
		ms["charset"] = "utf8"
	}
	db, err := gorm.Open(ms["datatype"], ms["username"]+":"+ms["password"]+"@tcp("+ms["hostname"]+")/"+ms["database"]+"?charset="+ms["charset"])
	if err != nil {
		panic(fmt.Sprintf("\x1b[31;20m[ERROR] %s\x1b[0m\n", err.Error()))
	}
	//全局禁用表复数
	db.SingularTable(true)
	DBS = append(DBS, name)
	DBC[name] = db
}
