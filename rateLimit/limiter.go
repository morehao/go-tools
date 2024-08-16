package rateLimit

import (
	"context"
	"errors"
	"github.com/go-redis/redis_rate/v10"
	"golang.org/x/time/rate"
	"time"
)

type Limiter interface {
	Allow(ctx context.Context, key string) (bool, error) // 是否允许请求
	Wait(ctx context.Context, key string) error          // 等待请求
}

func NewLimiter(opts ...Option) (Limiter, error) {
	cfg := &Config{
		Rate:            1,           // 默认每秒一个请求
		Burst:           1,           // 默认容量为1
		Period:          time.Second, // 默认时间窗口为1秒
		CleanupInterval: time.Minute, // 默认清理间隔为1分钟
	}
	for _, opt := range opts {
		opt(cfg)
	}
	switch cfg.Mode {
	case ModeTimeRate:
		limiter := &timeRateLimiter{
			limiterMap:      make(map[string]*rate.Limiter),
			lastAccessedMap: make(map[string]time.Time),
			burst:           cfg.Burst,
			period:          cfg.Period,
			cleanupInterval: cfg.CleanupInterval,
		}
		go limiter.cleanupLoop()
		return limiter, nil
	case ModeRedis:
		if cfg.RedisClient == nil {
			return nil, errors.New("redis client is nil")
		}
		return &redisLimiter{
			limiter: redis_rate.NewLimiter(cfg.RedisClient),
			rate:    cfg.Rate,
			burst:   cfg.Burst,
			period:  cfg.Period,
		}, nil
	default:
		return nil, errors.New("unsupported limiter mode")
	}
}
