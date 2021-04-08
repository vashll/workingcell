package workingcell

import (
	"os"
	"os/signal"
	"syscall"
	"workincell/common"
	"workincell/network"
)

func Start() {
	tcpcell := network.NewTcpCell("")
	tcpcell.StartServe("addr")
}

func stop() {
	close(common.StopChan)
}

func WaitExit(fns ...func()) {
	signal.Notify(common.SysStopChan, os.Interrupt, os.Kill, syscall.SIGTERM)
	select {
	case <-common.SysStopChan:
		stop()
	}
	for _, f := range fns {
		f()
	}
}