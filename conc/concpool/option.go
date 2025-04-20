package concpool

import (
	"context"
)

// Option 是一个函数类型，用于设置 pool 的选项
type Option func(q *pool)

// WithContext 设置 pool 的上下文
func WithContext(ctx context.Context) Option {
	return func(q *pool) {
		q.ctx = ctx
	}
}

// WithErrorHandler 设置 pool 的错误处理函数
func WithErrorHandler(handler func(err error)) Option {
	return func(q *pool) {
		q.onErr = handler
	}
}
