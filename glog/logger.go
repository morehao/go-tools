package glog

import (
	"context"
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
	getLogger(opts ...Option) (Logger, error)
	Close()
}

type LoggerConfig struct {
	LoggerType LoggerType `yaml:"logger_type"`
	Service    string     `yaml:"service"`
	Level      Level      `yaml:"level"`
	Dir        string     `yaml:"dir"`
	Stdout     bool       `yaml:"stdout"`
	ExtraKeys  []string   `yaml:"extra_keys"`
}

func NewLogger(cfg *LoggerConfig, opts ...Option) error {
	var logger Logger
	switch cfg.LoggerType {
	case LoggerTypeZap:
		l, err := newZapLogger(cfg, opts...)
		if err != nil {
			return err
		}
		logger = l
	default:
		l, err := newZapLogger(cfg, opts...)
		if err != nil {
			return err
		}
		logger = l
	}
	logInstance = instance{
		Logger:  logger,
		logType: cfg.LoggerType,
	}
	return nil
}

// newZapLogger 初始化zapLogger
func newZapLogger(cfg *LoggerConfig, opts ...Option) (Logger, error) {
	logger, err := getZapLogger(cfg)
	if err != nil {
		return nil, err
	}
	optCfg := &optConfig{}
	for _, opt := range opts {
		opt.apply(optCfg)
	}
	// AddCallerSkip(3) 跳过三层调用，使得日志输出正确的业务文件名和函数
	// logger = logger.WithOptions(zap.AddCallerSkip(3))
	return &zapLogger{
		logger: logger.WithOptions(optCfg.zapOpts...),
		cfg:    cfg,
	}, nil
}
