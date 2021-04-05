package log

import (
	"fmt"
	"runtime"
)

func LogInfo(format string, args ...interface{}) {
	fmt.Sprintf(format, args)
}

func LogError(format string, args ...interface{}) {
	fmt.Sprintf(format, args)
}

func LogStack() {
	buf := make([]byte, 1<<12)
	LogError(string(buf[:runtime.Stack(buf, false)]))
}
