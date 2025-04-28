package glog

import "sync"

var lock sync.RWMutex

func init() {
	// 初始化默认logger
	logger, err := getDefaultLogger()
	if err != nil {
		panic(err)
	}
	loggerInst := &loggerInstance{Logger: logger}
	lock.Lock()
	defaultLogger = loggerInst
	lock.Unlock()
}
