package glog

import (
	"context"
)

// Field 日志字段类型
type Field struct {
	Key   string
	Value interface{}
}

// Hook 钩子函数类型
type Hook func(ctx context.Context, level Level, msg string, fields ...Field)

// hooks 全局钩子函数列表
var hooks []Hook

// AddHook 添加钩子函数
func AddHook(hook Hook) {
	hooks = append(hooks, hook)
}

// executeHooks 执行所有钩子函数
func executeHooks(ctx context.Context, level Level, msg string, fields ...Field) {
	for _, hook := range hooks {
		func() {
			defer func() {
				if r := recover(); r != nil {
					// 钩子函数中的panic不应该影响日志记录
				}
			}()
			hook(ctx, level, msg, fields...)
		}()
	}
}

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
		logger: logger,
		cfg:    cfg,
	}, nil
}
