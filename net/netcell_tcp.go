package net

import (
	"net"
	"os"
	. "workincell/log"
)

type tcpNetCell struct {
	dataHandler IDataHandler
	stop bool
}

func(r *tcpNetCell) SetDataHandler(hanlder IDataHandler) {
	r.dataHandler = hanlder
}

func(r *tcpNetCell) StartServe(addr string) {
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		LogError("tcpcell start listen on addr:%s failed:%s", addr, err)
		os.Exit(1)
		return
	}
	for !r.stop {
		conn, err := listener.Accept()
		if err != nil {
			LogError("tcp cell accept fail :%s", err)
			break
		} else {
			msgCell := newTcpMsgCell(conn, r.dataHandler)
			msgCell.Run()
		}
	}
}