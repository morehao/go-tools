package glog

import (
	"context"
)

type Logger interface {
	Debug(ctx context.Context, args ...any)
	Debugf(ctx context.Context, format string, args ...any)
	Debugw(ctx context.Context, msg string, keysAndValues ...any)
	Info(ctx context.Context, args ...any)
	Infof(ctx context.Context, format string, args ...any)
	Infow(ctx context.Context, msg string, keysAndValues ...any)
	Warn(ctx context.Context, args ...any)
	Warnf(ctx context.Context, format string, args ...any)
	Warnw(ctx context.Context, msg string, keysAndValues ...any)
	Error(ctx context.Context, args ...any)
	Errorf(ctx context.Context, format string, args ...any)
	Errorw(ctx context.Context, msg string, keysAndValues ...any)
	Panic(ctx context.Context, args ...any)
	Panicf(ctx context.Context, format string, args ...any)
	Panicw(ctx context.Context, msg string, keysAndValues ...any)
	Fatal(ctx context.Context, args ...any)
	Fatalf(ctx context.Context, format string, args ...any)
	Fatalw(ctx context.Context, msg string, keysAndValues ...any)
	getLogger(opts ...Option) (Logger, error)
	Close()
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
	logInstance = &instance{
		Logger: logger,
	}
	return nil
}

// newZapLogger 初始化zapLogger
func newZapLogger(cfg *LoggerConfig, opts ...Option) (Logger, error) {
	optCfg := &optConfig{}
	for _, opt := range opts {
		opt.apply(optCfg)
	}
	logger, err := getZapLogger(cfg, optCfg)
	if err != nil {
		return nil, err
	}

	return &zapLogger{
		logger: logger.WithOptions(optCfg.zapOpts...),
		cfg:    cfg,
	}, nil
}
