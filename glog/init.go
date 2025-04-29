package glog

func init() {
	// 初始化默认logger
	logger, err := getDefaultLogger()
	if err != nil {
		panic(err)
	}
	loggerInst := &loggerInstance{Logger: logger}
	defaultLoggerInstance = loggerInst
}
