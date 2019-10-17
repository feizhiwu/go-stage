package common

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

var DB *gorm.DB

type DBInfo struct {
	Datatype string `yaml:"datatype"`
	Hostname string `yaml:"hostname"`
	Database string `yaml:"database"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
}

type Database struct {
	Debug   DBInfo
	Test    DBInfo
	Release DBInfo
}

func (d *Database) GetConnect() {
	database := d.getDatabase()
	DB, _ = gorm.Open(database.Datatype, database.Username+":"+database.Password+"@/"+database.Database+"?charset=utf8&parseTime=True&loc=Local")
}

//获取数据库配置
func (d *Database) getDatabase() DBInfo {
	dir, _ := os.Getwd()
	filePath := path.Join(dir, "/app/config/database.yml")
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
