package glog

import (
	"context"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Debugw(ctx context.Context, msg string, keysAndValues ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Infow(ctx context.Context, msg string, keysAndValues ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Warnw(ctx context.Context, msg string, keysAndValues ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Errorw(ctx context.Context, msg string, keysAndValues ...interface{})
	Panic(ctx context.Context, args ...interface{})
	Panicf(ctx context.Context, format string, args ...interface{})
	Panicw(ctx context.Context, msg string, keysAndValues ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
	Fatalw(ctx context.Context, msg string, keysAndValues ...interface{})
	Close()
}

type LoggerConfig struct {
	ServiceName string   `yaml:"service_name"`
	Level       Level    `yaml:"level"`
	LogDir      string   `yaml:"log_dir"`
	InConsole   bool     `yaml:"in_console"`
	ExtraKeys   []string `yaml:"extra_keys"`
}

// InitZapLogger 初始化zapLogger
func InitZapLogger(cfg *LoggerConfig) error {
	logger, err := newZapLogger(cfg)
	if err != nil {
		return err
	}
	// AddCallerSkip(3) 跳过三层调用，使得日志输出正确的业务文件名和函数
	logger = logger.WithOptions(zap.AddCallerSkip(3))
	logInstance = &zapLogger{
		logger: logger,
		cfg:    cfg,
	}
	return nil
}
