package network

import "workincell/common"

type WorkConfig struct {
	WorkType       int32
	MaxMsgQueueLen int32
	PoolSize       int32
}

type WorkerBuilder struct {
	worker IWorker
	cfg    *WorkConfig
}

func (r *WorkerBuilder) build() IWorker {
	if r.cfg.WorkType == common.WorkTypeUnique {
		if r.worker == nil {
			r.worker = newUniqueWorker(r.cfg.MaxMsgQueueLen)
		}
		return r.worker
	}
	if r.cfg.WorkType == common.WorkTypePool {
		if r.worker == nil {
			r.worker = newPoolWorker(r.cfg.PoolSize, r.cfg.MaxMsgQueueLen)
		}
		return r.worker
	}
	return newReactorWorker(r.cfg.MaxMsgQueueLen)
}

//func getDefault
