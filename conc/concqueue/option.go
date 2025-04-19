package concqueue

import (
	"context"
	"time"
)

// Option 是一个函数类型，用于设置 queue 的选项
type Option func(q *queue)

// ErrorHandler 是一个函数类型，用于处理任务执行过程中发生的错误
type ErrorHandler func(error)

// WithContext 设置 queue 的上下文
func WithContext(ctx context.Context) Option {
	return func(q *queue) {
		q.ctx = ctx
	}
}

// WithSubmitTimeout 设置 Submit 方法的超时时间
func WithSubmitTimeout(timeout time.Duration) Option {
	return func(q *queue) {
		q.submitTimeout = timeout
	}
}

// WithShutdownTimeout 设置 Shutdown 方法的超时时间
func WithShutdownTimeout(timeout time.Duration) Option {
	return func(q *queue) {
		q.shutdownTimeout = timeout
	}
}

// WithErrorHandler 设置 ErrorHandler
func WithErrorHandler(handler ErrorHandler) Option {
	return func(q *queue) {
		q.errorHandler = handler
	}
}

// PanicHandler 是一个函数类型，用于处理任务执行过程中发生的 panic
type PanicHandler func(interface{})

// WithPanicHandler 设置 PanicHandler
func WithPanicHandler(handler PanicHandler) Option {
	return func(q *queue) {
		q.panicHandler = handler
	}
}

// WithLogger 设置 Logger
func WithLogger(logger Logger) Option {
	return func(q *queue) {
		q.logger = logger
	}
}

// WithContextKeys 设置需要在日志中包含的 Context Key
func WithContextKeys(keys ...interface{}) Option {
	return func(q *queue) {
		q.contextKeys = keys
	}
}
