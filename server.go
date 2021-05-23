package workingcell

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"workincell/common"
	"workincell/network"
)

var ServerCfg *common.ServerConfig
var ExtWorkCfg *network.WorkConfig
var RpcWorkCfg *network.WorkConfig

func init() {
	ServerCfg = &common.ServerConfig{}
}

func SetServerAddr(ip string, extPort, rpcPort int) {
	ServerCfg.AddrIp = ip
	ServerCfg.ExtPort = extPort
	ServerCfg.RpcPort = rpcPort
}

func Start() {
	if ServerCfg.AddrIp == "" {
		return
	}
	if ServerCfg.ExtPort > 0 {
		addr := fmt.Sprintf("%s:%v", ServerCfg.AddrIp, ServerCfg.ExtPort)
		tcpCell := network.NewTcpCell(common.ConnTypeExt, addr, ExtWorkCfg)
		tcpCell.StartServe()
	}
	if ServerCfg.RpcPort > 0 {
		addr := fmt.Sprintf("%s:%v", ServerCfg.AddrIp, ServerCfg.RpcPort)
		tcpCell := network.NewTcpCell(common.ConnTypeRpc, addr, RpcWorkCfg)
		tcpCell.StartServe()
	}
}

func stop() {
	close(common.StopChan)
}

func WaitExit(fns ...func()) {
	signal.Notify(common.SysStopChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-common.SysStopChan:
		stop()
	}
	for _, f := range fns {
		f()
	}
}
