package net

import (
	"net"
	"workincell/common"
	"workincell/log"
)

type tcpMsgCell struct {
	conn net.Conn
	dataHandler IDataHandler
	stop bool
}

func newTcpMsgCell(conn net.Conn, hanlder IDataHandler) *tcpMsgCell {
	msgCell := &tcpMsgCell{}
	msgCell.conn = conn
	msgCell.dataHandler = hanlder
	return msgCell
}

func(r *tcpMsgCell) Run() {
	if r.dataHandler == nil {
		log.LogError("tcp msg cell run fail, data handler is nil")
		return
	}
	common.Go(func() {
			r.read()
	})
	common.Go(func() {
			r.write()
	})
}

func(r *tcpMsgCell) read() {
	for !r.stop {
		err, data := r.dataHandler.ReadData(r.conn)
		if err != nil {
			log.LogError("tcp msg cell read data fail :%v", err)
			break
		}
		msg := DefaultDataToMsg(data)
	}

}

func(r *tcpMsgCell) write() {

}