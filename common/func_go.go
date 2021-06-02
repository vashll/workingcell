package common

import (
	"os"
	"workincell/log"
	"sync/atomic"
)

var SysStopChan = make(chan os.Signal)
var StopChan = make(chan struct{})

var GoCount int32

func Go(fn func()) {
	atomic.AddInt32(&GoCount, 1)
	go func() {
		defer func() {
			atomic.AddInt32(&GoCount, -1)
			if err := recover(); err != nil {
				log.LogStack()
			}
		}()
		fn()
	}()
}

func GoWithCallBack(fn func(), cfn func()) {
	atomic.AddInt32(&GoCount, 1)
	go func() {
		defer func() {
			atomic.AddInt32(&GoCount, -1)
			if err := recover(); err != nil {
				log.LogStack()
			}
			cfn()
		}()
		fn()
	}()
}

func GetGoCount() int32 {
	return atomic.LoadInt32(&GoCount)
}
