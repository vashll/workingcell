package network

import (
	"workincell/common"
)

type WorkConfig struct {
	WorkType       int32
	MaxMsgQueueLen int32
	PoolSize       int32
}

type WorkerBuilder struct {
	worker IWorker
	cfg    *WorkConfig
}

func newWorkerBuilder(cfg *WorkConfig) *WorkerBuilder {
	builder := &WorkerBuilder{}
	builder.cfg = cfg
	return builder
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
	// WorkTypeReactor
	return newReactorWorker(r.cfg.MaxMsgQueueLen)
}

func getDefaultExtConfig() *WorkConfig {
	return &WorkConfig{
		WorkType:       common.WorkTypeReactor,
		MaxMsgQueueLen: 8,
		PoolSize:       0,
	}
}

func getDefaultRpcConfig() *WorkConfig {
	return &WorkConfig{
		WorkType:       common.WorkTypePool,
		MaxMsgQueueLen: 128,
		PoolSize:       4,
	}
}
