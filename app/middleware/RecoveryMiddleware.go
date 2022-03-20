package middleware

import (
	"fmt"
	"github.com/feizhiwu/gs/albedo"
	"github.com/gin-gonic/gin"
	"log"
	"runtime"
	"stage/app/common"
	"strings"
)

func Recovery(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			rollback(c)
			message := fmt.Sprintf("%s", err)
			logger := fmt.Sprintf("%s", trace(message))
			if c.Request.Context().Value("dbLog") != nil {
				logger += "\n" + albedo.MakeString(c.Request.Context().Value("dbLog"))
			}
			logger += "\n----------------------------------------------------------------------"
			for _, v := range writer() {
				log.New(v, "", log.LstdFlags).Printf("%s", logger)
			}

			common.NewDisplay(c).Show(err)
		}
	}()
	c.Next()
}

func trace(message string) string {
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
