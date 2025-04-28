package glog

import (
	"context"
)

type loggerInstance struct {
	Logger
}

// defaultLogger 默认的日志实例
var defaultLogger *loggerInstance

// moduleLoggerConfigMap 存储模块级别的logger
var moduleLoggerConfigMap = make(map[string]*ModuleLoggerConfig)

// InitLogger 初始化日志系统
func InitLogger(config *LogConfig, opts ...Option) error {
	lock.Lock()
	defer lock.Unlock()
	// 初始化模块级别的logger
	for module, cfg := range config.Modules {
		// 设置模块配置的 service 和 module 字段
		cfg.service = config.Service
		cfg.module = module
		moduleLoggerConfigMap[module] = cfg
	}

	// 设置默认logger
	defaultModuleConfig, ok := moduleLoggerConfigMap[defaultModuleName]
	if !ok {
		cfg := getDefaultModuleLoggerConfig()
		defaultModuleConfig = cfg
	}
	logger, err := newZapLogger(defaultModuleConfig, opts...)
	if err != nil {
		return err
	}
	defaultLogger = &loggerInstance{Logger: logger}

	return nil
}

func GetModuleLogger(moduleName string, opts ...Option) (Logger, error) {
	lock.RLock()
	defer lock.RUnlock()
	moduleLogConfig, ok := moduleLoggerConfigMap[moduleName]
	if !ok {
		defaultModuleCfg := getDefaultModuleLoggerConfig()
		defaultModuleCfg = defaultModuleCfg.ResetModule(moduleName)
		return newZapLogger(defaultModuleCfg, opts...)
	}
	logger, err := newZapLogger(moduleLogConfig, opts...)
	if err != nil {
		return nil, err
	}
	return &loggerInstance{Logger: logger}, nil
}

// getLoggerFromCtx 从Context中获取logger，如果没有则返回默认logger
func getLoggerFromCtx(ctx context.Context) Logger {
	logger, ok := ctx.Value(KeyLogger).(Logger)
	if ok {
		if logger == nil {
			return defaultLogger
		}
		return logger
	}
	return defaultLogger
}

// 以下函数使用Context中的logger，如果没有则使用默认logger

func Debug(ctx context.Context, args ...any) {
	defaultLogger.Debug(ctx, args...)
}

func Debugf(ctx context.Context, format string, kvs ...any) {
	defaultLogger.Debugf(ctx, format, kvs...)
}

func Debugw(ctx context.Context, msg string, kvs ...any) {
	defaultLogger.Debugw(ctx, msg, kvs...)
}

func Info(ctx context.Context, args ...any) {
	defaultLogger.Info(ctx, args...)
}

func Infof(ctx context.Context, format string, kvs ...any) {
	defaultLogger.Infof(ctx, format, kvs...)
}

func Infow(ctx context.Context, msg string, kvs ...any) {
	defaultLogger.Infow(ctx, msg, kvs...)
}

func Warn(ctx context.Context, args ...any) {
	defaultLogger.Warn(ctx, args...)
}

func Warnf(ctx context.Context, format string, kvs ...any) {
	defaultLogger.Warnf(ctx, format, kvs...)
}

func Warnw(ctx context.Context, msg string, kvs ...any) {
	defaultLogger.Warnw(ctx, msg, kvs...)
}

func Error(ctx context.Context, args ...any) {
	defaultLogger.Error(ctx, args...)
}

func Errorf(ctx context.Context, format string, kvs ...any) {
	defaultLogger.Errorf(ctx, format, kvs...)
}

func Errorw(ctx context.Context, msg string, kvs ...any) {
	defaultLogger.Errorw(ctx, msg, kvs...)
}

func Panic(ctx context.Context, args ...any) {
	defaultLogger.Panic(ctx, args...)
}

func Panicf(ctx context.Context, format string, kvs ...any) {
	defaultLogger.Panicf(ctx, format, kvs...)
}

func Panicw(ctx context.Context, msg string, kvs ...any) {
	defaultLogger.Panicw(ctx, msg, kvs...)
}

func Fatal(ctx context.Context, args ...any) {
	defaultLogger.Fatal(ctx, args...)
}

func Fatalf(ctx context.Context, format string, kvs ...any) {
	defaultLogger.Fatalf(ctx, format, kvs...)
}

func Fatalw(ctx context.Context, msg string, kvs ...any) {
	defaultLogger.Fatalw(ctx, msg, kvs...)
}

func Name(ctx context.Context) string {
	return defaultLogger.Name()
}

// Close 关闭所有logger
func Close() {
	defaultLogger.Close()
}
