package glog

import (
	"context"
)

type instance struct {
	Logger
}

var logInstance *instance

func Debug(ctx context.Context, args ...any) {
	logInstance.Debug(ctx, args...)
}

func Debugf(ctx context.Context, format string, args ...any) {
	logInstance.Debugf(ctx, format, args...)
}

func Debugw(ctx context.Context, msg string, keysAndValues ...any) {
	logInstance.Debugw(ctx, msg, keysAndValues...)
}

func Info(ctx context.Context, args ...any) {
	logInstance.Info(ctx, args...)
}

func Infof(ctx context.Context, format string, args ...any) {
	logInstance.Infof(ctx, format, args...)
}

func Infow(ctx context.Context, msg string, keysAndValues ...any) {
	logInstance.Infow(ctx, msg, keysAndValues...)
}

func Warn(ctx context.Context, args ...any) {
	logInstance.Warn(ctx, args...)
}

func Warnf(ctx context.Context, format string, args ...any) {
	logInstance.Warnf(ctx, format, args...)
}

func Warnw(ctx context.Context, msg string, keysAndValues ...any) {
	logInstance.Warnw(ctx, msg, keysAndValues...)
}

func Error(ctx context.Context, args ...any) {
	logInstance.Error(ctx, args...)
}
func Errorf(ctx context.Context, format string, args ...any) {
	logInstance.Errorf(ctx, format, args...)
}

func Errorw(ctx context.Context, msg string, keysAndValues ...any) {
	logInstance.Errorw(ctx, msg, keysAndValues...)
}

func Panic(ctx context.Context, args ...any) {
	logInstance.Panic(ctx, args...)
}

func Panicf(ctx context.Context, format string, args ...any) {
	logInstance.Panicf(ctx, format, args...)
}

func Fatal(ctx context.Context, args ...any) {
	logInstance.Fatal(ctx, args...)
}

func Fatalf(ctx context.Context, format string, args ...any) {
	logInstance.Fatalf(ctx, format, args...)
}

func Fatalw(ctx context.Context, msg string, keysAndValues ...any) {
	logInstance.Fatalw(ctx, msg, keysAndValues...)
}

func GetLogger(opts ...Option) (Logger, error) {
	if logInstance == nil {
		cfg := getDefaultLoggerConfig()
		logger, err := newZapLogger(cfg, opts...)
		if err != nil {
			return nil, err
		}
		logInstance = &instance{Logger: logger}
		return logger, nil
	}
	return logInstance.Logger.getLogger(opts...)
}

func Close() {
	logInstance.Close()
}
