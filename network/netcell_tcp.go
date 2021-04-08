package network

import (
	"net"
	"os"
	"workincell/common"
	. "workincell/log"
)

type tcpNetCell struct {
	listener   net.Listener
	dataReader IDataReader
	stop       bool
}

func (r *tcpNetCell) SetDataReader(reader IDataReader) {
	r.dataReader = reader
}

func (r *tcpNetCell) StartServe(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		LogError("tcpcell start listen on addr:%s failed:%s", addr, err)
		os.Exit(1)
		return
	}
	r.listener = listener
	common.Go(func() {
		for !r.stop {
			conn, err := listener.Accept()
			if err != nil {
				LogError("tcp cell accept fail :%s", err)
				break
			} else {
				msgCell := newTcpMsgCell(conn, r.dataReader)
				msgCell.Run()
			}
		}
	})
}

func NewTcpCell(nettyp string) *tcpNetCell {
	netcell := &tcpNetCell{}
	netcell.dataReader = &DefaultDataReader{}
	return netcell
}
