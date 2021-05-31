package log

import (
	"fmt"
	"runtime"
)

const (
	LogLevelInfo  = 0
	LogLevelDebug = 1
	LogLevelWarn  = 3
	LogLevelError = 4
)

var logLevel = LogLevelInfo

func SetLogLevel(level int) {
	logLevel = level
}

func LogInfo(format string, args ...interface{}) {
	if logLevel >= LogLevelInfo {
		fmt.Println(fmt.Sprintf(format, args...))
	}
}

func LogError(format string, args ...interface{}) {
	if logLevel >= LogLevelError {
		fmt.Println(fmt.Sprintf(format, args...))
	}
}

func LogWarn(format string, args ...interface{}) {
	if logLevel >= LogLevelWarn {
		fmt.Println(fmt.Sprintf(format, args...))
	}
}

func LogStack() {
	buf := make([]byte, 1<<12)
	s := string(buf[:runtime.Stack(buf, false)])
	fmt.Println(s)
}
