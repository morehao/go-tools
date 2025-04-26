package glog

import (
	"context"
)

// defaultLogger 默认的日志实例
var defaultLogger Logger

// moduleLoggers 存储模块级别的logger
var moduleLoggers = make(map[string]Logger)

// SetDefaultLogger 设置默认logger
func SetDefaultLogger(logger Logger) {
	defaultLogger = logger
}

// GetLogger 从Context中获取logger，如果没有则返回默认logger
func GetLogger(ctx context.Context) Logger {
	if logger, ok := ctx.Value(loggerKey).(Logger); ok {
		return logger
	}
	return defaultLogger
}

// WithLogger 将logger注入到Context中
func WithLogger(ctx context.Context, logger Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}

// loggerKey 用于在Context中存储logger的key
type loggerKeyType struct{}

var loggerKey = loggerKeyType{}

// 以下函数使用Context中的logger，如果没有则使用默认logger

func Debug(ctx context.Context, kvs ...any) {
	GetLogger(ctx).Debug(ctx, kvs...)
}

func Debugf(ctx context.Context, format string, kvs ...any) {
	GetLogger(ctx).Debugf(ctx, format, kvs...)
}

func Debugw(ctx context.Context, msg string, kvs ...any) {
	GetLogger(ctx).Debugw(ctx, msg, kvs...)
}

func Info(ctx context.Context, kvs ...any) {
	GetLogger(ctx).Info(ctx, kvs...)
}

func Infof(ctx context.Context, format string, kvs ...any) {
	GetLogger(ctx).Infof(ctx, format, kvs...)
}

func Infow(ctx context.Context, msg string, kvs ...any) {
	GetLogger(ctx).Infow(ctx, msg, kvs...)
}

func Warn(ctx context.Context, kvs ...any) {
	GetLogger(ctx).Warn(ctx, kvs...)
}

func Warnf(ctx context.Context, format string, kvs ...any) {
	GetLogger(ctx).Warnf(ctx, format, kvs...)
}

func Warnw(ctx context.Context, msg string, kvs ...any) {
	GetLogger(ctx).Warnw(ctx, msg, kvs...)
}

func Error(ctx context.Context, kvs ...any) {
	GetLogger(ctx).Error(ctx, kvs...)
}

func Errorf(ctx context.Context, format string, kvs ...any) {
	GetLogger(ctx).Errorf(ctx, format, kvs...)
}

func Errorw(ctx context.Context, msg string, kvs ...any) {
	GetLogger(ctx).Errorw(ctx, msg, kvs...)
}

func Panic(ctx context.Context, kvs ...any) {
	GetLogger(ctx).Panic(ctx, kvs...)
}

func Panicf(ctx context.Context, format string, kvs ...any) {
	GetLogger(ctx).Panicf(ctx, format, kvs...)
}

func Panicw(ctx context.Context, msg string, kvs ...any) {
	GetLogger(ctx).Panicw(ctx, msg, kvs...)
}

func Fatal(ctx context.Context, kvs ...any) {
	GetLogger(ctx).Fatal(ctx, kvs...)
}

func Fatalf(ctx context.Context, format string, kvs ...any) {
	GetLogger(ctx).Fatalf(ctx, format, kvs...)
}

func Fatalw(ctx context.Context, msg string, kvs ...any) {
	GetLogger(ctx).Fatalw(ctx, msg, kvs...)
}

// GetModuleLogger 获取指定模块的logger
func GetModuleLogger(module string) Logger {
	lock.RLock()
	defer lock.RUnlock()
	if logger, ok := moduleLoggers[module]; ok {
		return logger
	}
	return defaultLogger
}

// gotLogger 获取指定模块的logger，如果不存在则返回默认logger
func gotLogger(module string) Logger {
	if logger, ok := moduleLoggers[module]; ok {
		return logger
	}
	return defaultLogger
}

// Close 关闭所有logger
func Close() {
	defaultLogger.Close()
}
