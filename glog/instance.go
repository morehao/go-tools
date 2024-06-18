package glog

import "context"

var logInstance Logger

func Debug(ctx context.Context, args ...interface{}) {
	logInstance.Debug(ctx, args...)
}

func Debugf(ctx context.Context, format string, args ...interface{}) {
	logInstance.Debugf(ctx, format, args...)
}

func Info(ctx context.Context, args ...interface{}) {
	logInstance.Info(ctx, args...)
}

func Infof(ctx context.Context, format string, args ...interface{}) {
	logInstance.Infof(ctx, format, args...)
}

func Warn(ctx context.Context, args ...interface{}) {
	logInstance.Warn(ctx, args...)
}

func Warnf(ctx context.Context, format string, args ...interface{}) {
	logInstance.Warnf(ctx, format, args...)
}

func Error(ctx context.Context, args ...interface{}) {
	logInstance.Error(ctx, args...)
}
func Errorf(ctx context.Context, format string, args ...interface{}) {
	logInstance.Errorf(ctx, format, args...)
}

func Panic(ctx context.Context, args ...interface{}) {
	logInstance.Panic(ctx, args...)
}

func Panicf(ctx context.Context, format string, args ...interface{}) {
	logInstance.Panicf(ctx, format, args...)
}

func Fatal(ctx context.Context, args ...interface{}) {
	logInstance.Fatal(ctx, args...)
}

func Fatalf(ctx context.Context, format string, args ...interface{}) {
	logInstance.Fatalf(ctx, format, args...)
}
