package concqueue

import (
	"context"
	"fmt"
	"io"
	"log"
)

type Logger interface {
	Errorf(ctx context.Context, format string, args ...any)
}

// defaultLogger 是默认的 Logger
type defaultLogger struct {
	logger      *log.Logger
	contextKeys []any
}

// newDefaultLogger 创建一个默认的 Logger
func newDefaultLogger(out io.Writer, prefix string, flag int, contextKeys []any) Logger {
	return &defaultLogger{
		logger:      log.New(out, prefix, flag),
		contextKeys: contextKeys,
	}
}

// Errorf 实现了 Logger 接口的 Errorf 方法
func (l *defaultLogger) Errorf(ctx context.Context, format string, args ...any) {
	// 从 Context 中获取需要的信息
	contextInfo := ""
	for _, key := range l.contextKeys {
		if value := ctx.Value(key); value != nil {
			contextInfo += fmt.Sprintf("%v=%v ", key, value)
		}
	}

	l.logger.Printf(contextInfo+format, args...)
}
