package glog

func init() {
	// 初始化默认logger
	cfg := &LoggerConfig{
		Module: "default",
		Level:  InfoLevel,
		Type:   WriterConsole,
		Dir:    "./log",
	}
	var err error
	defaultLogger, err = newZapLogger(cfg)
	if err != nil {
		panic(err)
	}
}

// Init 初始化日志系统
func Init(config *ServiceConfig, opts ...Option) error {
	// 初始化模块级别的logger
	for module, cfg := range config.Modules {
		logger, err := newZapLogger(cfg, opts...)
		if err != nil {
			return err
		}
		moduleLoggers[module] = logger
	}

	// 设置默认logger
	if defaultCfg, ok := config.Modules["default"]; ok {
		logger, err := newZapLogger(defaultCfg, opts...)
		if err != nil {
			return err
		}
		defaultLogger = logger
	}

	return nil
}
