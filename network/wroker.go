package network

import (
	"sync"
	"sync/atomic"
	"workincell/common"
	"workincell/log"
)

var UniqueWorker *WorkCell
var uniqueLock sync.Locker
var MsgDispachter *WorkCell

type WorkCell struct {
	msgChan chan *Message
	state   int32
}

func NewWorker(wtye int32, maxMsgQueueLen int32) *WorkCell {
	if wtye == 1 {
		//唯一worker模式
		if UniqueWorker == nil {
			uniqueLock.Lock()
			if UniqueWorker == nil {
				UniqueWorker = createWorker(maxMsgQueueLen)
			}
			uniqueLock.Unlock()
		}
		return UniqueWorker
	} else if wtye == 2 {
		//工作池模式 pending

	} else if wtye == 3 {
		//连接独享worker
		return createWorker(maxMsgQueueLen)
	}
	return nil
}

func createWorker(maxMsgQueueLen int32) *WorkCell {
	wc := &WorkCell{}
	if maxMsgQueueLen < 1 {
		log.LogError("worker msg queue len less than 1.")
		return nil
	}
	wc.msgChan = make(chan *Message, maxMsgQueueLen)
	wc.startWork()
	return wc
}

func (r *WorkCell) IsRunning() bool {
	return r.state == 1
}

func (r *WorkCell) Stop() {
	r.state = 1
}

func (r *WorkCell) startWork() {
	if 0 == atomic.LoadInt32(&r.state) {
		return
	}
	for r.IsRunning() {
		select {
		case msg, ok := <-r.msgChan:
			if ok {
				onMsg(msg)
			}
		case <-common.StopChan:
			r.Stop()
			break
		}
	}
}

func (r *WorkCell) PushMsg(msg *Message) {
	r.msgChan <- msg
}

func onMsg(msg *Message) {

}
