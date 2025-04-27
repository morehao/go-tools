package glog

import (
	"context"
)

// defaultLogger 默认的日志实例
var defaultLogger Logger

// moduleLoggers 存储模块级别的logger
var moduleLoggers = make(map[string]Logger)

// InitLogger 初始化日志系统
func InitLogger(config *LogConfig, opts ...Option) error {
	lock.Lock()
	defer lock.Unlock()
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
	defaultLogger = moduleLoggers[defaultModuleName]
	if defaultLogger == nil {
		logger, err := getDefaultLogger()
		if err != nil {
			return err
		}
		defaultLogger = logger
		moduleLoggers[defaultModuleName] = logger
	}

	return nil
}

func GetModuleLogger(moduleName string) (Logger, error) {
	lock.RLock()
	defer lock.RUnlock()
	logger, ok := moduleLoggers[moduleName]
	if !ok {
		moduleCfg := getDefaultModuleLoggerConfig()
		moduleCfg = moduleCfg.ResetModule(moduleName)
		return newZapLogger(moduleCfg)
	}
	return logger, nil
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

func Debug(ctx context.Context, kvs ...any) {
	getLoggerFromCtx(ctx).Debug(ctx, kvs...)
}

func Debugf(ctx context.Context, format string, kvs ...any) {
	getLoggerFromCtx(ctx).Debugf(ctx, format, kvs...)
}

func Debugw(ctx context.Context, msg string, kvs ...any) {
	getLoggerFromCtx(ctx).Debugw(ctx, msg, kvs...)
}

func Info(ctx context.Context, kvs ...any) {
	getLoggerFromCtx(ctx).Info(ctx, kvs...)
}

func Infof(ctx context.Context, format string, kvs ...any) {
	getLoggerFromCtx(ctx).Infof(ctx, format, kvs...)
}

func Infow(ctx context.Context, msg string, kvs ...any) {
	getLoggerFromCtx(ctx).Infow(ctx, msg, kvs...)
}

func Warn(ctx context.Context, kvs ...any) {
	getLoggerFromCtx(ctx).Warn(ctx, kvs...)
}

func Warnf(ctx context.Context, format string, kvs ...any) {
	getLoggerFromCtx(ctx).Warnf(ctx, format, kvs...)
}

func Warnw(ctx context.Context, msg string, kvs ...any) {
	getLoggerFromCtx(ctx).Warnw(ctx, msg, kvs...)
}

func Error(ctx context.Context, kvs ...any) {
	getLoggerFromCtx(ctx).Error(ctx, kvs...)
}

func Errorf(ctx context.Context, format string, kvs ...any) {
	getLoggerFromCtx(ctx).Errorf(ctx, format, kvs...)
}

func Errorw(ctx context.Context, msg string, kvs ...any) {
	getLoggerFromCtx(ctx).Errorw(ctx, msg, kvs...)
}

func Panic(ctx context.Context, kvs ...any) {
	getLoggerFromCtx(ctx).Panic(ctx, kvs...)
}

func Panicf(ctx context.Context, format string, kvs ...any) {
	getLoggerFromCtx(ctx).Panicf(ctx, format, kvs...)
}

func Panicw(ctx context.Context, msg string, kvs ...any) {
	getLoggerFromCtx(ctx).Panicw(ctx, msg, kvs...)
}

func Fatal(ctx context.Context, kvs ...any) {
	getLoggerFromCtx(ctx).Fatal(ctx, kvs...)
}

func Fatalf(ctx context.Context, format string, kvs ...any) {
	getLoggerFromCtx(ctx).Fatalf(ctx, format, kvs...)
}

func Fatalw(ctx context.Context, msg string, kvs ...any) {
	getLoggerFromCtx(ctx).Fatalw(ctx, msg, kvs...)
}

// Close 关闭所有logger
func Close() {
	defaultLogger.Close()
}
