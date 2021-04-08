package common

import (
	"os"
	. "workincell/log"
	)

var SysStopChan = make(chan os.Signal)
var StopChan = make(chan struct{})

func Go(fn func()) {
	//id := atomic.AddUint32(&goid, 1)
	//c := atomic.AddInt32(&gocount, 1)
	go func() {
		defer func() {
			if err := recover(); err != nil {
				LogStack()
			}
		}()
		fn()
	}()
}