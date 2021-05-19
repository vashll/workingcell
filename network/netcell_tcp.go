package network

import (
	"net"
	"time"
	"workincell/common"
	"workincell/log"
)

type tcpNetCell struct {
	listener      net.Listener
	dataReader    IDataReader
	stop          int32
	addr          string
	connType      int32 //连接类型(client,rpc)
	workerBuilder *WorkerBuilder
}

func (r *tcpNetCell) SetDataReader(reader IDataReader) {
	r.dataReader = reader
}

func (r *tcpNetCell) IsStop() bool {
	return r.stop == 1
}

func (r *tcpNetCell) StartServe() {
	var addr string
	if r.connType == common.ConnTypeClient {
		addr = common.ServerCfg.AddrIp + ":" + common.ServerCfg.TcpPort
	} else if r.connType == common.ConnTypeRpc {
		addr = common.ServerCfg.AddrIp + ":" + common.ServerCfg.RpcPort
	} else {
		log.LogError("connection type is error, type:%v", r.connType)
		return
	}
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.LogError("tcpcell start listen on addr:%s failed:%s", addr, err)
		return
	}
	r.listener = listener
	r.stop = 0
	common.Go(func() {
		select {
		case <-common.StopChan:
			if r.listener != nil {
				r.stop = 1
				r.listener.Close()
			}
		}
	})
	common.Go(func() {
		for !r.IsStop() {
			conn, err := listener.Accept()
			if err != nil {
				log.LogError("tcp cell accept fail :%s", err)
				break
			} else {
				msgCell := newTcpMsgCell(conn, r.dataReader, 1, 1, 1)
				msgCell.Run()
			}
		}
	})
}

func NewTcpCell(connType int32, addr string, workCfg *WorkConfig) *tcpNetCell {
	netcell := &tcpNetCell{}
	netcell.dataReader = &DefaultDataReader{}
	netcell.connType = connType
	netcell.addr = addr

	return netcell
}

func TcpConnect(addr string, connType int32, reader IDataReader) {
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		log.LogError("connect to addr:%s fail err:%s", addr, err.Error())
		return
	}
	msgCell := newTcpMsgCell(conn, reader, 1, 1, 1)
	msgCell.Run()
}
