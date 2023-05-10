package scene

import (
	"context"
	"github.com/gin-gonic/gin"
)

func Access(c *gin.Context, data map[string]interface{}) {
	//TODO 访问日志
}

func Error(ctx context.Context, error string) {
	//TODO 错误日志
}
