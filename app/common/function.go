package common

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"stage/config/conf"
	"strconv"
	"strings"
	"time"
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

func MakeString(str interface{}) string {
	switch str.(type) {
	case string:
		return str.(string)
	case int, uint, int8, uint8, int16, uint16, int32, uint32, int64, uint64:
		return strconv.Itoa(MakeInt(str))
	case float32:
		return strconv.FormatFloat(float64(str.(float32)), 'f', -1, 32)
	case float64:
		return strconv.FormatFloat(str.(float64), 'f', -1, 32)
	default:
		return ""
	}
}

func GetOffset(page interface{}, limit uint) uint {
	num := MakeUint(page) - 1
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

func EncryptPass(pass string) string {
	salt := conf.GetValue("salt").(string)
	sum := md5.Sum([]byte(pass + salt))
	return fmt.Sprintf("%x", sum)
}

func Time() string {
	return strconv.FormatInt(time.Now().Unix(), 10)
}

//获取格式日期
func Date(timestamp interface{}, format string) string {
	sec := int64(MakeUint(timestamp))
	date := time.Unix(sec, 0).Format(format)
	return date
}

//获取全格式日期
func TimeFormat(timestamp interface{}) string {
	return Date(timestamp, "2006-01-02 15:04:05")
}

//日期转换时间戳
func Timestamp(date string) string {
	var layout string
	s := strings.Split(date, " ")
	if len(s) > 1 {
		layout = "2006-01-02 15:04:05"
	} else {
		s := strings.Split(date, "-")
		if len(s) > 2 {
			layout = "2006-01-02"
		} else {
			layout = "2006-01"
		}
	}
	location, _ := time.LoadLocation("Local")
	timestamp, _ := time.ParseInLocation(layout, date, location)
	return strconv.FormatInt(timestamp.Unix(), 10)
}

func InArray(search interface{}, array interface{}) bool {
	switch n := search.(type) {
	case uint:
		for _, v := range array.([]uint) {
			if n == v {
				return true
			}
		}
	case int:
		for _, v := range array.([]int) {
			if n == v {
				return true
			}
		}
	case string:
		for _, v := range array.([]string) {
			if n == v {
				return true
			}
		}
	default:
		return false
	}
	return false
}

func Implode(array interface{}, sep string) string {
	return strings.Replace(strings.Trim(fmt.Sprint(array), "[]"), " ", sep, -1)
}
