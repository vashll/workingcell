package common

import (
	"encoding/json"
	"io"
	"os"
	"workincell/log"
)

type ServerConfig struct {
	AddrIp  string `json:"addr_ip"`
	ExtPort int    `json:"ext_port"`
	RpcPort int    `json:"rpc_port"`
}

const (
	WorkTypeUnique  = 1 //唯一协程模式
	WorkTypePool    = 2 //协程池模式
	WorkTypeReactor = 3 //Reactor模式
)

//连接类型
const (
	ConnTypeExt = 1 //外部连接
	ConnTypeRpc = 2 //RPC连接
)

//var ServerCfg *ServerConfig

//func GenServerConfig() *ServerConfig {
//	if configPath == "" {
//		ServerCfg = &ServerConfig{
//			AddrIp:  "127.0.0.1",
//			TcpPort: "7777",
//			RpcPort: "8888",
//		}
//	} else {
//		readServerConfig()
//	}
//	return ServerCfg
//}

func readServerConfig(path string) (serverCfg *ServerConfig) {
	if path == "" {
		log.LogError("config path is nil.")
		return
	}
	f, err := os.Open(path)
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
	serverCfg = &ServerConfig{}
	err = json.Unmarshal(data, serverCfg)
	if err != nil {
		serverCfg = nil
		log.LogError("config unmaishal fail, ", err.Error())
		return
	}
	return
}
