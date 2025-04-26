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

func Debug(ctx context.Context, args ...any) {
	GetLogger(ctx).Debug(ctx, args...)
}

func Debugf(ctx context.Context, format string, args ...any) {
	GetLogger(ctx).Debugf(ctx, format, args...)
}

func Debugw(ctx context.Context, msg string, keysAndValues ...any) {
	GetLogger(ctx).Debugw(ctx, msg, keysAndValues...)
}

func Info(ctx context.Context, args ...any) {
	GetLogger(ctx).Info(ctx, args...)
}

func Infof(ctx context.Context, format string, args ...any) {
	GetLogger(ctx).Infof(ctx, format, args...)
}

func Infow(ctx context.Context, msg string, keysAndValues ...any) {
	GetLogger(ctx).Infow(ctx, msg, keysAndValues...)
}

func Warn(ctx context.Context, args ...any) {
	GetLogger(ctx).Warn(ctx, args...)
}

func Warnf(ctx context.Context, format string, args ...any) {
	GetLogger(ctx).Warnf(ctx, format, args...)
}

func Warnw(ctx context.Context, msg string, keysAndValues ...any) {
	GetLogger(ctx).Warnw(ctx, msg, keysAndValues...)
}

func Error(ctx context.Context, args ...any) {
	GetLogger(ctx).Error(ctx, args...)
}

func Errorf(ctx context.Context, format string, args ...any) {
	GetLogger(ctx).Errorf(ctx, format, args...)
}

func Errorw(ctx context.Context, msg string, keysAndValues ...any) {
	GetLogger(ctx).Errorw(ctx, msg, keysAndValues...)
}

func Panic(ctx context.Context, args ...any) {
	GetLogger(ctx).Panic(ctx, args...)
}

func Panicf(ctx context.Context, format string, args ...any) {
	GetLogger(ctx).Panicf(ctx, format, args...)
}

func Fatal(ctx context.Context, args ...any) {
	GetLogger(ctx).Fatal(ctx, args...)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	GetLogger(ctx).Fatalf(ctx, format, args...)
}

func Fatalw(ctx context.Context, msg string, keysAndValues ...any) {
	GetLogger(ctx).Fatalw(ctx, msg, keysAndValues...)
}

// GetModuleLogger 获取指定模块的logger
func GetModuleLogger(module string) Logger {
	// 这里可以扩展为支持模块级别的logger
	return defaultLogger
}

// Close 关闭所有logger
func Close() {
	defaultLogger.Close()
}

// gotLogger 获取指定模块的logger，如果不存在则返回默认logger
func gotLogger(module string) Logger {
	if logger, ok := moduleLoggers[module]; ok {
		return logger
	}
	return defaultLogger
}
