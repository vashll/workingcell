package network

import (
	"net"
	"sync/atomic"
	"workincell/common"
	"workincell/log"
)

type TcpMsgCell struct {
	conn       net.Conn
	writeCh    chan *Message
	dataReader IDataReader
	stop       int32
	worker     *WorkCell
	agentData  interface{}
}

func newTcpMsgCell(conn net.Conn, reader IDataReader) *TcpMsgCell {
	msgCell := &TcpMsgCell{}
	msgCell.conn = conn
	msgCell.dataReader = reader
	msgCell.writeCh = make(chan *Message, 8)
	msgCell.worker = NewWorker(3, 16)
	return msgCell
}

func (r *TcpMsgCell) IsStop() bool {
	return r.stop == 1
}

//Set the message cell to run state.
func (r *TcpMsgCell) SetStart() {
	r.stop = 1
}

//Start a message cell.
func (r *TcpMsgCell) Run() {
	if r.dataReader == nil {
		log.LogError("tcp msg cell run fail, data handler is nil")
		return
	}
	r.SetStart()
	common.Go(func() {
		r.read()
	})
	common.Go(func() {
		r.write()
	})
}

//Stop a message cell, include read and write goroutine.
func (r *TcpMsgCell) Stop() {
	atomic.CompareAndSwapInt32(&r.stop, 1, 0)
	if r.conn != nil {
		r.conn.Close()
	}
}

//Read data form connection by data reader
func (r *TcpMsgCell) read() {
	for !r.IsStop() {
		err, msg := r.dataReader.ReadData(r.conn)
		if err != nil {
			log.LogError("tcp msg cell read data fail :%v", err)
			break
		}
		if r.worker != nil {
			cellMsg := &CellMessage{
				Msg:  msg,
				User: r.agentData,
			}
			r.worker.PushMsg(cellMsg)
		}
	}
}

// Read message form input channel and write
func (r *TcpMsgCell) write() {
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

// Send message to client
func (r *TcpMsgCell) Send(msg *Message) {
	if msg == nil || r.IsStop() {
		return
	}
	select {
	case r.writeCh <- msg:
	default:
		log.LogError("message cell send crowed, canceled send.")
	}
}

func (r *TcpMsgCell) SetAgentData(ad interface{}) {
	r.agentData = ad
}

func (r *TcpMsgCell) GetAgentDate() interface{} {
	return r.agentData
}
