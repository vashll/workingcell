package common

import (
	. "workincell/log"
	)

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