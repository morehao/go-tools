package glog

import "sync"

var lock sync.RWMutex

func init() {
	// 初始化默认logger
	cfg := &ModuleLoggerConfig{
		module:  "default",
		Level:   InfoLevel,
		Writer:  WriterConsole,
		Dir:     "./log",
		service: "app",
	}
	var err error
	defaultLogger, err = newZapLogger(cfg)
	if err != nil {
		panic(err)
	}
	lock.Lock()
	moduleLoggers[defaultModuleName] = defaultLogger
	lock.Unlock()
}

// Init 初始化日志系统
func Init(config *LogConfig, opts ...Option) error {
	// 初始化模块级别的logger
	for module, cfg := range config.Modules {
		// 设置模块配置的 service 和 module 字段
		cfg.service = config.Service
		cfg.module = module
		logger, err := newZapLogger(cfg, opts...)
		if err != nil {
			return err
		}
		lock.Lock()
		moduleLoggers[module] = logger
		lock.Unlock()
	}

	// 设置默认logger
	lock.Lock()
	defaultLogger = moduleLoggers[defaultModuleName]
	lock.Unlock()
	if defaultLogger == nil {
		// 如果没有默认logger，创建一个
		cfg := getDefaultLoggerConfig()
		cfg.service = config.Service
		cfg.module = defaultModuleName
		logger, err := newZapLogger(cfg, opts...)
		if err != nil {
			return err
		}
		defaultLogger = logger
	}

	return nil
}
