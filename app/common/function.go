package common

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"strconv"
)

func MakeUint(num interface{}) uint {
	switch num.(type) {
	case float32:
		return uint(num.(float32))
	case float64:
		return uint(num.(float64))
	case string:
		i, _ := strconv.Atoi(num.(string))
		return uint(i)
	default:
		return 0
	}
}

func MakeInt(num interface{}) int {
	return int(MakeUint(num))
}

func MakeJson(data map[string]interface{}) []byte {
	bytes, _ := json.Marshal(data)
	return bytes
}

func GetOffset(page interface{}, limit int) int {
	num := MakeInt(page) - 1
	if num < 0 {
		num = 0
	}
	return num * limit
}

func GetData(c *gin.Context) map[string]interface{} {
	data, _ := c.Get("param")
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
