package concqueue

import (
	"context"
)

// Option 是一个函数类型，用于设置 queue 的选项
type Option func(q *queue)

// WithContext 设置 queue 的上下文
func WithContext(ctx context.Context) Option {
	return func(q *queue) {
		q.ctx = ctx
	}
}
