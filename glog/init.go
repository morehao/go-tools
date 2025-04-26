package glog

func init() {
	// 初始化默认logger
	cfg := &LoggerConfig{
		module:  "default",
		Level:   InfoLevel,
		Type:    WriterConsole,
		Dir:     "./log",
		service: "app",
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
		// 设置模块配置的 service 和 module 字段
		cfg.service = config.Service
		cfg.module = module
		logger, err := newZapLogger(cfg, opts...)
		if err != nil {
			return err
		}
		moduleLoggers[module] = logger
	}

	// 设置默认logger
	defaultLogger = moduleLoggers["default"]
	if defaultLogger == nil {
		defaultLogger = moduleLoggers["app"]
	}
	if defaultLogger == nil {
		// 如果没有默认logger，创建一个
		cfg := getDefaultLoggerConfig()
		cfg.service = config.Service
		cfg.module = "default"
		logger, err := newZapLogger(cfg, opts...)
		if err != nil {
			return err
		}
		defaultLogger = logger
	}

	return nil
}
