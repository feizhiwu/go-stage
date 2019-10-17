package controller

import (
	"toutGin/app/common"
)

var data map[string]interface{}
var display *common.Display

func init() {
	data = make(map[string]interface{})
}
