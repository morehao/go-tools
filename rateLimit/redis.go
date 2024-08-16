package rateLimit

import (
	"context"
	"github.com/go-redis/redis_rate/v10"
	"time"
)

type redisLimiter struct {
	limiter *redis_rate.Limiter
	rate    int           // 每秒允许的最大请求数
	period  time.Duration // 限制周期
}

func (l *redisLimiter) Allow(ctx context.Context, key string) (bool, error) {
	res, err := l.limiter.Allow(ctx, key, redis_rate.Limit{
		Rate:   l.rate,
		Period: l.period,
		Burst:  l.rate,
	})
	if err != nil {
		return false, err
	}
	return res.Allowed > 0, nil
}

func (l *redisLimiter) Wait(ctx context.Context, key string) error {
	for {
		res, err := l.limiter.Allow(ctx, key, redis_rate.Limit{
			Rate:   l.rate,
			Period: l.period,
			Burst:  l.rate,
		})
		if err != nil {
			return err
		}
		if res.Allowed > 0 {
			break
		}

		// 如果不允许，等待 RetryAfter 指定的时间再重试
		select {
		case <-time.After(res.RetryAfter):
			// 继续循环检查
		case <-ctx.Done():
			// 如果上下文被取消，返回上下文的错误
			return ctx.Err()
		}
	}
	return nil
}
