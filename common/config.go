package common

import (
	"encoding/json"
	"io"
	"os"
	"workincell/log"
)

type ServerConfig struct {
	WorkType       int32 `json:"work_type"`
	MaxMsgQueueLen int32 `json:"max_msg_queue_len"`
	PoolSize       int32 `json:"pool_size"`
}

const (
	WorkTypeUnique  = 1 //唯一协程模式
	WorkTypePool    = 2 //协程池模式
	WorkTypeReactor = 3 //Reactor模式
)

var configPath string
var ServerCfg *ServerConfig

func SetServerConfigPath(path string) {
	configPath = path
}

func GenServerConfig() *ServerConfig {
	if configPath == "" {
		ServerCfg = &ServerConfig{}
		//默认模式 Reactor
		ServerCfg.WorkType = WorkTypeReactor
		ServerCfg.MaxMsgQueueLen = 8
	} else {
		readServerConfig()
	}
	return ServerCfg
}

func readServerConfig() {
	if configPath == "" {
		log.LogError("config path is nil.")
		return
	}
	f, err := os.Open(configPath)
	if err != nil {
		log.LogError("open config file fail, ", err.Error())
		return
	}
	defer func() {
		f.Close()
	}()
	data, err := io.ReadAll(f)
	if err != nil {
		log.LogError("read config file fail, ", err.Error())
		return
	}
	ServerCfg = &ServerConfig{}
	err = json.Unmarshal(data, ServerCfg)
	if err != nil {
		ServerCfg = nil
		log.LogError("config unmaishal fail, ", err.Error())
		return
	}
}
