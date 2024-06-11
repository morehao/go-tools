package glog

import (
	"context"
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
	ServiceName string `yaml:"service_name"`
	Level       Level  `yaml:"level"`
	LogDir      string `yaml:"log_dir"`
	InConsole   bool   `yaml:"in_console"`
}

// NewZapLogger 创建一个新的 Logger 实例
func NewZapLogger(cfg *LoggerConfig) (Logger, error) {
	logger, err := newZapLogger(cfg)
	if err != nil {
		return nil, err
	}
	return &zapLogger{
		logger: logger,
		cfg:    cfg,
	}, nil
}
