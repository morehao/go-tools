package glog

import "sync"

var lock sync.RWMutex

func init() {
	// 初始化默认logger
	logger, err := getDefaultLogger()
	if err != nil {
		panic(err)
	}
	lock.Lock()
	moduleLoggers[defaultModuleName] = logger
	defaultLogger = logger
	lock.Unlock()
}
