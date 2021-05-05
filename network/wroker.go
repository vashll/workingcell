package network

import (
	"sync"
	"sync/atomic"
	"workincell/common"
	"workincell/log"
)

const (
	WorkStateRunning = 0
	WorkStateStopped = 1
)

var Worker *WorkCell
var uniqueLock sync.Locker

type WorkCell struct {
	workTyp  int32 //工作模式
	msgChan  chan *CellMessage
	state    int32 // 0:running, 1:stopped
	poolChan chan int32
	poolSize int32
}

func NewWorker() *WorkCell {
	cfg := common.ServerCfg
	if cfg.WorkType == common.WorkTypeUnique {
		//唯一worker模式
		if Worker == nil {
			uniqueLock.Lock()
			if Worker == nil {
				Worker = createWorker(cfg.WorkType, cfg.MaxMsgQueueLen)
			}
			uniqueLock.Unlock()
		}
		return Worker
	} else if cfg.WorkType == common.WorkTypePool {
		if Worker == nil {
			uniqueLock.Lock()
			if Worker == nil {
				Worker = createWorker(cfg.WorkType, cfg.MaxMsgQueueLen)
				Worker.poolSize = cfg.PoolSize
			}
			uniqueLock.Unlock()
		}
	} else if cfg.WorkType == common.WorkTypeReactor {
		//连接独享worker
		return createWorker(cfg.WorkType, cfg.MaxMsgQueueLen)
	}
	return nil
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
