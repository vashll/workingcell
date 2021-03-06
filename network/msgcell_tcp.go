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
	worker     IWorker
	agentData  interface{}
}

func newTcpMsgCell(conn net.Conn, reader IDataReader, worker IWorker) *TcpMsgCell {
	msgCell := &TcpMsgCell{}
	msgCell.conn = conn
	msgCell.dataReader = reader
	msgCell.writeCh = make(chan *Message, 8)
	msgCell.worker = worker
	return msgCell
}

func (r *TcpMsgCell) IsStop() bool {
	return r.stop == 1
}

//Set the message cell to run state.
func (r *TcpMsgCell) SetStart() {
	r.stop = 0
	if r.worker != nil {
		r.worker.StartWork()
	}
}

//Start a message cell.
func (r *TcpMsgCell) Run() {
	if r.dataReader == nil {
		log.LogError("tcp msg cell run fail, data handler is nil")
		return
	}
	r.SetStart()
	r.read()
	r.write()
	log.LogInfo("== go count:%v", common.GetGoCount())
}

//Stop a message cell, include read and write goroutine.
func (r *TcpMsgCell) Stop() {
	atomic.CompareAndSwapInt32(&r.stop, 1, 0)
	if r.conn != nil {
		r.conn.Close()
	}
	if r.writeCh != nil {
		close(r.writeCh)
	}
	if r.worker != nil {
		r.worker.Stop()
	}
}

//Read data form connection by data reader
func (r *TcpMsgCell) read() {
	common.Go(func() {
		for !r.IsStop() {
			msg, err := r.dataReader.ReadData(r.conn)
			if err != nil {
				log.LogError("tcp msg cell read data fail :%v", err)
				break
			}
			r.PushMsg(msg)
		}
	})
}

// Read message form input channel and write
func (r *TcpMsgCell) write() {
	common.Go(func() {
		var msg *Message
		for !r.IsStop() {
			select {
			case msg = <-r.writeCh:
			case <-common.StopChan:
				r.Stop()
				continue
			}
			if msg == nil {
				continue
			}
			if r.conn != nil {
				r.conn.Write(r.dataReader.MsgToData(msg))
			}
		}
	})
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

//push msg to worker
func (r *TcpMsgCell) PushMsg(msg *Message) {
	if r.worker != nil {
		cellMsg := &CellMessage{
			Msg:     msg,
			User:    r.agentData,
			MsgCell: r,
		}
		r.worker.PushMsg(cellMsg)
	}
}

func (r *TcpMsgCell) SetAgentData(ad interface{}) {
	r.agentData = ad
}

func (r *TcpMsgCell) GetAgentData() interface{} {
	return r.agentData
}
