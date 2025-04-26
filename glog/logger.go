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
	Debug(ctx context.Context, kvs ...any)
	Debugf(ctx context.Context, format string, kvs ...any)
	Debugw(ctx context.Context, msg string, kvs ...any)
	Info(ctx context.Context, kvs ...any)
	Infof(ctx context.Context, format string, kvs ...any)
	Infow(ctx context.Context, msg string, kvs ...any)
	Warn(ctx context.Context, kvs ...any)
	Warnf(ctx context.Context, format string, kvs ...any)
	Warnw(ctx context.Context, msg string, kvs ...any)
	Error(ctx context.Context, kvs ...any)
	Errorf(ctx context.Context, format string, kvs ...any)
	Errorw(ctx context.Context, msg string, kvs ...any)
	Panic(ctx context.Context, kvs ...any)
	Panicf(ctx context.Context, format string, kvs ...any)
	Panicw(ctx context.Context, msg string, kvs ...any)
	Fatal(ctx context.Context, kvs ...any)
	Fatalf(ctx context.Context, format string, kvs ...any)
	Fatalw(ctx context.Context, msg string, kvs ...any)
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

// convertKvsToFields 将 kvs ...any 转换为 []Field
func convertKvsToFields(kvs ...any) []Field {
	fields := make([]Field, 0, len(kvs)/2)
	for i := 0; i < len(kvs); i += 2 {
		if i+1 >= len(kvs) {
			break
		}
		key, ok := kvs[i].(string)
		if !ok {
			continue
		}
		fields = append(fields, Field{Key: key, Value: kvs[i+1]})
	}
	return fields
}
