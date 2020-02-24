package config

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var DB *gorm.DB

func GetValue(key string) interface{} {
	dir, _ := os.Getwd()
	filePath := path.Join(dir, "/config/config.yml")
	fileData, _ := ioutil.ReadFile(filePath)
	var config map[interface{}]interface{}
	yaml.Unmarshal(fileData, &config)
	if gin.Mode() == gin.ReleaseMode {
		config = config["release"].(map[interface{}]interface{})
	} else {
		config = config["test"].(map[interface{}]interface{})
	}
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

type DBInfo struct {
	Datatype string `yaml:"datatype"`
	Hostname string `yaml:"hostname"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Prefix   string `yaml:"prefix"`
}

type Database struct {
	Debug   DBInfo
	Test    DBInfo
	Release DBInfo
}

func (d *Database) GetConnect() {
	database := d.GetInfo()
	DB, _ = gorm.Open(database.Datatype, database.Username+":"+database.Password+"@/"+database.Database+"?charset=utf8&parseTime=True&loc=Local")
}

//获取数据库配置
func (d *Database) GetInfo() DBInfo {
	dir, _ := os.Getwd()
	filePath := path.Join(dir, "/config/database.yml")
	fileData, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(fileData, &d)
	if gin.Mode() == gin.DebugMode {
		return d.Debug
	} else if gin.Mode() == gin.TestMode {
		return d.Test
	} else {
		return d.Release
	}
}

type Message struct {
	Msg map[int]string
}

//根据status返回文字说明
func (m *Message) GetMessage(status int) string {
	var filePath string
	dir, _ := os.Getwd()
	filePath = path.Join(dir, "/config/message.yml")
	fileData, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(fileData, &m)
	res := m.Msg[status]
	if res == "" {
		return m.GetMessage(11000)
	}
	return res
}
