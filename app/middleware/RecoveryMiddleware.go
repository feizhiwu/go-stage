package middleware

import (
	"context"
	"fmt"
	"github.com/feizhiwu/gs/albedo"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
	"stage/app/plugin"
	"stage/app/plugin/scene"
	"stage/config/conf"
	"strings"
	"time"
)

func Recovery(c *gin.Context) {
	//设置访问时间
	defer func() {
		if err := recover(); err != nil {
			ctx := c.Request.Context()
			t := ctx.Value("accessTime").(time.Time)
			level := "OK"
			if time.Since(t) > time.Second*1 {
				level = "SLOW"
			}
			logger := fmt.Sprintf("%s[%s] %s %s action:%s %v in %v", c.ClientIP(), level, c.Request.Method,
				c.Request.RequestURI, c.GetHeader("action"), plugin.GetParams(c), time.Since(t))
			logger += "\n" + fmt.Sprintf("%s", trace(err))
			if ctx.Value("query") != nil {
				logger += "\n" + albedo.MakeString(ctx.Value("query"))
			}
			logger += "\n----------------------------------------------------------------------"
			for _, v := range writer() {
				log.New(v, "", log.LstdFlags).Printf("%s", logger)
			}
			if val, ok := err.(int); ok {
				scene.Error(ctx, fmt.Sprintf("%d：%s", val, conf.Message(val)))
			} else {
				scene.Error(ctx, fmt.Sprintf("%s", err))
			}

			scene.Action(c).Show(err)
		}
	}()
	ctx := c.Request.Context()
	ctx = context.WithValue(c.Request.Context(), "accessTime", time.Now())
	c.Request = c.Request.WithContext(ctx)
	c.Next()
}

func trace(err interface{}) string {
	message := fmt.Sprintf("%s", err)
	var pcs [32]uintptr
	n := runtime.Callers(3, pcs[:])
	var str strings.Builder
	str.WriteString(message + "\nTraceback:")
	for _, pc := range pcs[:n] {
		fn := runtime.FuncForPC(pc)
		file, line := fn.FileLine(pc)
		str.WriteString(fmt.Sprintf("\n\t%s:%d", file, line))
	}
	return str.String()
}
