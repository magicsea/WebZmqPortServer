package config

import (
	json "encoding/json"
	"io/ioutil"
)

type WebConfig struct {
	PostURL     string
	OnException string //异常返回的结果[nil:不返回，#ERRMSG#返回错误消息,#ERRSTACK#返回错误堆栈] 示例：have error:#ERRMSG#
	MQRemote    string
	FormNames   []string
}

var (
	ServerConfig WebConfig
)

func ReadConfigFromFile(configFileName string) error {
	ServerConfig = WebConfig{}

	data, err := ioutil.ReadFile(configFileName)
	if err != nil {
		return err
	}

	dataBytes := []byte(data)
	err = json.Unmarshal(dataBytes, &ServerConfig)
	return err
}
