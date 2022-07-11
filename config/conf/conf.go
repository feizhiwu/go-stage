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
	"time"
)

var (
	DBC  map[string]*gorm.DB
	DBS  []string
	File *file
)

type file struct {
	date    string
	pack    string
	Config  string
	Message string
	LogDir  string
	Logger  string
}

func GetFile() *file {
	if File == nil || File.date != time.Now().Format("02") {
		File = new(file)
		File.date = time.Now().Format("02")
		if gin.Mode() == gin.DebugMode {
			dir, _ := os.Getwd()
			File.Config = path.Join(dir, "config/config.yml")
			File.Message = path.Join(dir, "config/message.yml")
			File.LogDir = path.Join(dir, "data/runtime/log/"+time.Now().Format("200601"))
		} else {
			File.pack = "/stage"
			File.Config = path.Join(os.Getenv("GOPATH"), "src", File.pack, "config/config.yml")
			File.Message = path.Join(os.Getenv("GOPATH"), "src", File.pack, "config/message.yml")
			File.LogDir = path.Join(os.Getenv("GOBIN"), "data", File.pack, "runtime/log/"+time.Now().Format("200601"))
		}
		File.Logger = path.Join(File.LogDir, File.date+"-"+gin.Mode()+".log")
	}
	return File
}

// Config 获取配置信息
func Config(key string) interface{} {
	filePath := GetFile().Config
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

// Message 获取状态文本
func Message(status int) string {
	var msg map[int]string
	var filePath string
	filePath = GetFile().Message
	fileData, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(fileData, &msg)
	return msg[status]
}

// ConnectDB 连接数据库
func ConnectDB() {
	DBC = make(map[string]*gorm.DB)
	dbConf := Config("db")
	for k, v := range dbConf.(map[interface{}]interface{}) {
		if _, ok := v.(map[interface{}]interface{}); !ok {
			connectDB("db", dbConf.(map[interface{}]interface{}))
			break
		} else {
			connectDB(albedo.MakeString(k), v.(map[interface{}]interface{}))
		}
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
	//保存连接信息，键名与配置信息同名
	DBS = append(DBS, name)
	DBC[name] = db
}
