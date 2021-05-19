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
	workTyp        int32 //工作模式
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
				break
			}
		}
	})
}

func (r *UniqueWorkCell) OnNewMsg(msg *CellMessage) {
	processMsg(msg.User, msg.MsgCell, msg.Msg)
}

func (r *UniqueWorkCell) PushMsg(msg *CellMessage) {
	//maybe drop msg when crowed
	r.msgChan <- msg
}

//工作池模式worker
type PoolWorkCell struct {
	workTyp        int32 //工作模式
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
					break
				}
			}
		})
	}
}

func (r *PoolWorkCell) OnNewMsg(msg *CellMessage) {
	processMsg(msg.User, msg.MsgCell, msg.Msg)
}

func (r *PoolWorkCell) PushMsg(msg *CellMessage) {
	//maybe drop msg when crowed
	r.msgChan <- msg
}

//Reactor
type ReactorWorkCell struct {
	workTyp        int32 //工作模式
	msgChan        chan *CellMessage
	state          int32
	maxMsgQueueLen int32
	poolSize       int32
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
				break
			}
		}
	})
}

func (r *ReactorWorkCell) OnNewMsg(msg *CellMessage) {
	processMsg(msg.User, msg.MsgCell, msg.Msg)
}

func (r *ReactorWorkCell) PushMsg(msg *CellMessage) {
	//maybe drop msg when crowed
	r.msgChan <- msg
}

//=========================================================================
type WorkCell struct {
	workTyp  int32 //工作模式
	msgChan  chan *CellMessage
	state    int32 // 0:running, 1:stopped
	poolChan chan int32
	poolSize int32
}

func createWorker(typ int32, maxMsgQueueLen int32) *WorkCell {
	wc := &WorkCell{}
	if maxMsgQueueLen < 1 {
		log.LogError("worker msg queue len less than 1.")
		return nil
	}
	wc.msgChan = make(chan *CellMessage, maxMsgQueueLen)
	wc.startWork()
	return wc
}

func (r *WorkCell) IsRunning() bool {
	return r.state == WorkStateRunning
}

func (r *WorkCell) Stop() {
	r.state = WorkStateStopped
}

func (r *WorkCell) startWork() {
	if 0 == atomic.LoadInt32(&r.state) {
		return
	}
	common.Go(func() {
		for r.IsRunning() {
			select {
			case msg, ok := <-r.msgChan:
				if ok {
					r.onNewMsg(msg)
				}
			case <-common.StopChan:
				r.Stop()
				break
			}
		}
	})
}

func (r *WorkCell) PushMsg(msg *CellMessage) {
	r.msgChan <- msg
}

func (r *WorkCell) onNewMsg(msg *CellMessage) {
	if r.workTyp == common.WorkTypePool {
		r.poolChan <- 1
		common.GoWithCallBack(func() {
			processMsg(msg.User, msg.MsgCell, msg.Msg)
		}, func() {
			<-r.poolChan
		})
	} else {
		processMsg(msg.User, msg.MsgCell, msg.Msg)
	}
}

func processMsg(data interface{}, msgCell *TcpMsgCell, msg *Message) {

}
