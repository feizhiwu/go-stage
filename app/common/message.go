package common

import (
	"encoding/json"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

type Message struct {
	Msg string
}

//根据status返回文字说明
func (m *Message) GetMessage(status int) string {
	dir, _ := os.Getwd()
	filePath := path.Join(dir, "/app/config/message.yml")
	fileData, _ := ioutil.ReadFile(filePath)
	yaml.Unmarshal(fileData, &m)
	var message map[int]string
	json.Unmarshal([]byte(m.Msg), &message)
	res := message[status]
	if res == "" {
		return m.GetMessage(11000)
	}
	return res
}
