package network

import (
	"sync/atomic"
	"workincell/common"
	"workincell/log"
)

const (
	WorkStateRunning = 1
	WorkStateStopped = 2
)

type IWorker interface {
	IsRunning() bool
	Stop()
	StartWork()
	PushMsg(msg *CellMessage)
	OnNewMsg(msg *CellMessage)
}

//唯一worker
type UniqueWorkCell struct {
	// workTyp        int32 //工作模式
	msgChan        chan *CellMessage
	state          int32
	maxMsgQueueLen int32
}

func newUniqueWorker(maxLen int32) *UniqueWorkCell {
	w := &UniqueWorkCell{}
	if maxLen <= 0 {
		maxLen = 256
	}
	w.state = WorkStateStopped
	w.maxMsgQueueLen = maxLen
	w.msgChan = make(chan *CellMessage, maxLen)
	return w
}

func (r *UniqueWorkCell) IsRunning() bool {
	return r.state == WorkStateRunning
}

func (r *UniqueWorkCell) Stop() {
	r.state = WorkStateStopped
}

func (r *UniqueWorkCell) StartWork() {
	if WorkStateRunning == atomic.LoadInt32(&r.state) {
		return
	}
	r.state = WorkStateRunning
	common.Go(func() {
		for r.IsRunning() {
			select {
			case msg, ok := <-r.msgChan:
				if ok {
					r.OnNewMsg(msg)
				}
			case <-common.StopChan:
				r.Stop()
			}
		}
	})
}

func (r *UniqueWorkCell) OnNewMsg(msg *CellMessage) {
	processMsg(msg.User, msg.MsgCell, msg.Msg)
}

func (r *UniqueWorkCell) PushMsg(msg *CellMessage) {
	//drop msg when crowed
	select {
	case r.msgChan <- msg:
	default:
		log.LogInfo("message queue is full, drop message")
	}
}

//工作池模式worker
type PoolWorkCell struct {
	// workTyp        int32 //工作模式
	msgChan        chan *CellMessage
	state          int32
	maxMsgQueueLen int32
	poolSize       int32
}

func newPoolWorker(maxLen, poolSize int32) *PoolWorkCell {
	w := &PoolWorkCell{}
	if maxLen <= 0 {
		maxLen = 256
	}
	if poolSize <= 0 {
		poolSize = 8
	}
	w.state = WorkStateStopped
	w.maxMsgQueueLen = maxLen
	w.poolSize = poolSize
	w.msgChan = make(chan *CellMessage, maxLen)
	return w
}

func (r *PoolWorkCell) IsRunning() bool {
	return r.state == WorkStateRunning
}

func (r *PoolWorkCell) Stop() {
	r.state = WorkStateStopped
}

func (r *PoolWorkCell) StartWork() {
	if WorkStateRunning == atomic.LoadInt32(&r.state) {
		return
	}
	r.state = WorkStateRunning
	for i := 0; i < int(r.poolSize); i++ {
		common.Go(func() {
			for r.IsRunning() {
				select {
				case msg, ok := <-r.msgChan:
					if ok {
						r.OnNewMsg(msg)
					}
				case <-common.StopChan:
					r.Stop()
				}
			}
		})
	}
}

func (r *PoolWorkCell) OnNewMsg(msg *CellMessage) {
	processMsg(msg.User, msg.MsgCell, msg.Msg)
}

func (r *PoolWorkCell) PushMsg(msg *CellMessage) {
	//drop msg when crowed
	select {
	case r.msgChan <- msg:
	default:
		log.LogInfo("message queue is full, drop message")
	}
}

//Reactor
type ReactorWorkCell struct {
	// workTyp        int32 //工作模式
	msgChan        chan *CellMessage
	state          int32
	maxMsgQueueLen int32
}

func newReactorWorker(maxLen int32) *ReactorWorkCell {
	w := &ReactorWorkCell{}
	if maxLen <= 0 {
		maxLen = 8
	}
	w.state = WorkStateStopped
	w.maxMsgQueueLen = maxLen
	w.msgChan = make(chan *CellMessage, maxLen)
	return w
}

func (r *ReactorWorkCell) IsRunning() bool {
	return r.state == WorkStateRunning
}

func (r *ReactorWorkCell) Stop() {
	r.state = WorkStateStopped
}

func (r *ReactorWorkCell) StartWork() {
	r.state = WorkStateRunning
	common.Go(func() {
		for r.IsRunning() {
			select {
			case msg, ok := <-r.msgChan:
				if ok {
					r.OnNewMsg(msg)
				}
			case <-common.StopChan:
				r.Stop()
			}
		}
	})
}

func (r *ReactorWorkCell) OnNewMsg(msg *CellMessage) {
	processMsg(msg.User, msg.MsgCell, msg.Msg)
}

func (r *ReactorWorkCell) PushMsg(msg *CellMessage) {
	//drop msg when crowed
	select {
	case r.msgChan <- msg:
	default:
		log.LogInfo("message queue is full, drop message")
	}
}

func processMsg(data interface{}, msgCell *TcpMsgCell, msg *Message) {
}
