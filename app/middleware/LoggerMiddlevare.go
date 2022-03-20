package middleware

import (
	"context"
	"fmt"
	"github.com/feizhiwu/gs/albedo"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"log"
	"os"
	"path"
	"stage/app/common"
	"stage/config/conf"
	"strings"
	"time"
)

//*************************DB log start*************************

type dbLog struct {
	c    *gin.Context
	name string
	db   *gorm.DB
}

func newDBLog(c *gin.Context, name string, db *gorm.DB) {
	db.LogMode(true)
	db.SetLogger(&dbLog{
		c,
		name,
		db,
	})
}

func (l *dbLog) Print(values ...interface{}) {
	var dbLog string
	level := values[0]
	ctx := l.c.Request.Context()
	if level == "sql" {
		var (
			//source    = values[1]
			queryTime = values[2]
			sql       = values[3]
			param     = values[4]
		)
		if queryTime.(time.Duration) > time.Millisecond*300 {
			level = "SLOW"
		} else {
			level = "OK"
		}
		dbLog = fmt.Sprintf("[%s][%s] %s %s in %s", l.name, level, sql, param, queryTime)
	} else {
		dbLog = fmt.Sprintf("[%s]%s", l.name, values)
	}
	if ctx.Value("dbLog") != nil {
		dbLog = albedo.MakeString(ctx.Value("dbLog")) + "\n" + dbLog
	}
	ctx = context.WithValue(l.c.Request.Context(), "dbLog", dbLog)
	l.c.Request = l.c.Request.WithContext(ctx)
}

//*************************DB log end*************************
func begin(c *gin.Context) {
	var db *gorm.DB
	var ctx context.Context
	for k, v := range conf.DBC {
		db = v.Begin()
		ctx = context.WithValue(c.Request.Context(), k, db)
		c.Request = c.Request.WithContext(ctx)
		newDBLog(c, strings.ToUpper(k), db)
	}
}

func commit(c *gin.Context) {
	ctx := c.Request.Context()
	for _, v := range conf.DBS {
		ctx.Value(v).(*gorm.DB).Commit()
	}
}

func rollback(c *gin.Context) {
	ctx := c.Request.Context()
	for _, v := range conf.DBS {
		ctx.Value(v).(*gorm.DB).Rollback()
	}
}

// Logger 日志middleware
func Logger(c *gin.Context) {
	t := time.Now()
	begin(c)
	c.Next()
	commit(c)
	level := "OK"
	if time.Since(t) > time.Second*1 {
		level = "SLOW"
	}
	logger := fmt.Sprintf("%s[%s] %s %s action:%s %v in %v", c.ClientIP(), level, c.Request.Method,
		c.Request.RequestURI, c.GetHeader("action"), common.GetParams(c), time.Since(t))
	if c.Request.Context().Value("dbLog") != nil {
		logger += "\n" + albedo.MakeString(c.Request.Context().Value("dbLog"))
	}
	logger += "\n----------------------------------------------------------------------"
	for _, v := range writer() {
		log.New(v, "", log.LstdFlags).Printf("%s", logger)
	}
}

func writer() (writer []*os.File) {
	dir, _ := os.Getwd()
	logPath := path.Join(dir, "/data/runtime/log/"+time.Now().Format("200601"))
	//创建文件路径
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		os.MkdirAll(logPath, os.ModePerm)
	}
	var logFix string
	if gin.Mode() == gin.TestMode {
		logFix = "-test.log"
	} else {
		logFix = ".log"
	}
	filePath := logPath + "/" + time.Now().Format("02") + logFix
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.Create(filePath)
	}
	logFile, _ := os.OpenFile(filePath, os.O_WRONLY|os.O_APPEND, os.ModePerm)
	writer = append(writer, logFile)
	//debug模式，在控制台输出日志信息
	if gin.Mode() == gin.DebugMode {
		writer = append(writer, os.Stdout)
	}
	return
}
