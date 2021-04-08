package network

import (
	"net"
	"sync/atomic"
	"workincell/common"
	"workincell/log"
)

type tcpMsgCell struct {
	conn       net.Conn
	writeCh    chan *Message
	dataReader IDataReader
	stop       int32
}

func newTcpMsgCell(conn net.Conn, reader IDataReader) *tcpMsgCell {
	msgCell := &tcpMsgCell{}
	msgCell.conn = conn
	msgCell.dataReader = reader
	msgCell.writeCh = make(chan *Message, 16)
	return msgCell
}

func (r *tcpMsgCell) IsStop() bool {
	return r.stop == 1
}

func (r *tcpMsgCell) Start() {
	r.stop = 1
}

func (r *tcpMsgCell) Run() {
	if r.dataReader == nil {
		log.LogError("tcp msg cell run fail, data handler is nil")
		return
	}
	r.Start()
	common.Go(func() {
		r.read()
	})
	common.Go(func() {
		r.write()
	})
}

func (r *tcpMsgCell) Stop() {
	atomic.CompareAndSwapInt32(&r.stop, 1, 0)
	if r.conn != nil {
		r.conn.Close()
	}
}

func (r *tcpMsgCell) read() {
	for !r.IsStop() {
		err, msg := r.dataReader.ReadData(r.conn)
		if err != nil {
			log.LogError("tcp msg cell read data fail :%v", err)
			break
		}
		onMessage(msg)
	}
}

func onMessage(msg *Message) {

}

func (r *tcpMsgCell) write() {
	var msg *Message
	for !r.IsStop() {
		select {
		case msg = <-r.writeCh:
		case <-common.StopChan:
			r.Stop()
			break
		}
		if msg == nil {
			continue
		}
		if r.conn != nil {
			r.conn.Write(r.dataReader.MsgToData(msg))
		}
	}
}
