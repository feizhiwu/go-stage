package common

import (
	"crypto/md5"
	"fmt"
	"github.com/gin-gonic/gin"
	"reflect"
	"runtime"
	"stage/config/conf"
	"strings"
)

func GetParams(c *gin.Context) map[string]interface{} {
	data, _ := c.Get("params")
	return data.(map[string]interface{})
}

func CopyParams(val []string, data map[string]interface{}) map[string]interface{} {
	params := make(map[string]interface{})
	for _, v := range val {
		if data[v] != nil {
			params[v] = data[v]
		}
	}
	return params
}

func EncryptPass(pass string) string {
	salt := conf.Config("salt").(string)
	sum := md5.Sum([]byte(pass + salt))
	return fmt.Sprintf("%x", sum)
}

func InArray(n int, f func(int) bool) bool {
	for i := 0; i < n; i++ {
		if f(i) {
			return true
		}
	}
	return false
}

func GetFuncName(f interface{}) string {
	name := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	arr := strings.Split(name, ".")
	return strings.Split(arr[len(arr)-1], "-")[0]
}
