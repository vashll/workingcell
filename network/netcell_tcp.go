package network

import (
	"net"
	"os"
	"time"
	"workincell/common"
	"workincell/log"
)

type tcpNetCell struct {
	listener      net.Listener
	dataReader    IDataReader
	stop          int32
	listenAddr    string
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
	if r.connType != common.ConnTypeExt && r.connType != common.ConnTypeRpc {
		log.LogError("connection type is error, type:%v", r.connType)
		return
	}
	addr := r.listenAddr
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
				msgCell := newTcpMsgCell(conn, r.dataReader, r.workerBuilder.build())
				msgCell.Run()
			}
		}
	})
}

func NewTcpCell(connType int32, addr string, workCfg *WorkConfig) *tcpNetCell {
	netcell := &tcpNetCell{}
	netcell.dataReader = &DefaultDataReader{}
	netcell.connType = connType
	netcell.listenAddr = addr
	if workCfg == nil {
		if connType == common.ConnTypeExt {
			workCfg = getDefaultExtConfig()
		} else if connType == common.ConnTypeRpc {
			workCfg = getDefaultRpcConfig()
		} else {
			log.LogError("undefined connection type")
			os.Exit(1)
			return nil
		}
	}
	netcell.workerBuilder = newWorkerBuilder(workCfg)
	return netcell
}

func TcpConnect(addr string, connType int32, reader IDataReader, workCfg *WorkConfig) {
	if workCfg == nil {
		if connType == common.ConnTypeExt {
			workCfg = getDefaultExtConfig()
		} else if connType == common.ConnTypeRpc {
			workCfg = getDefaultRpcConfig()
		} else {
			log.LogError("undefined connection type")
			return
		}
	}
	conn, err := net.DialTimeout("tcp", addr, 2*time.Second)
	if err != nil {
		log.LogError("connect to addr:%s fail err:%s", addr, err.Error())
		return
	}
	wb := newWorkerBuilder(workCfg)
	msgCell := newTcpMsgCell(conn, reader, wb.build())
	msgCell.Run()
}
