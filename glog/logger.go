package glog

import (
	"context"
	"go.uber.org/zap"
)

type Logger interface {
	Debug(ctx context.Context, args ...interface{})
	Debugf(ctx context.Context, format string, args ...interface{})
	Info(ctx context.Context, args ...interface{})
	Infof(ctx context.Context, format string, args ...interface{})
	Warn(ctx context.Context, args ...interface{})
	Warnf(ctx context.Context, format string, args ...interface{})
	Error(ctx context.Context, args ...interface{})
	Errorf(ctx context.Context, format string, args ...interface{})
	Panic(ctx context.Context, args ...interface{})
	Panicf(ctx context.Context, format string, args ...interface{})
	Fatal(ctx context.Context, args ...interface{})
	Fatalf(ctx context.Context, format string, args ...interface{})
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
	// AddCallerSkip(2) 跳过两层调用，使得日志输出正确的业务文件名和函数
	logger = logger.WithOptions(zap.AddCallerSkip(3))
	logInstance = &zapLogger{
		logger: logger,
		cfg:    cfg,
	}
	return nil
}
