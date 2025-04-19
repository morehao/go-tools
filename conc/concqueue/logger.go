package concqueue

import "log"

// Logger 是一个接口，用于记录日志
type Logger interface {
	Printf(format string, v ...interface{})
}

// defaultLogger 是默认的 Logger
type defaultLogger struct{}

// Printf 实现了 Logger 接口的 Printf 方法
func (l *defaultLogger) Printf(format string, v ...interface{}) {
	log.Printf(format, v...)
}
