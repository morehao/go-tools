package concqueue

import "context"

type Option func(q *queue)

func WithContext(ctx context.Context) Option {
	return func(q *queue) {
		q.ctx = ctx
	}
}
