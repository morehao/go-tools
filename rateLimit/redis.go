package rateLimit

import (
	"context"
	"errors"
	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
	"sync"
	"sync/atomic"
	"time"
)

const pingInterval = time.Millisecond * 100

type redisLimiter struct {
	limiter        *redis_rate.Limiter
	client         *redis.Client
	rate           int           // 限制周期内允许的最大请求数
	burst          int           // 限制周期突发内允许的请求数
	period         time.Duration // 限制周期
	rescueLock     sync.Mutex
	redisAlive     uint32
	monitorStarted bool
	rescueLimiter  *timeRateLimiter
}

func (l *redisLimiter) Allow(ctx context.Context, key string) (bool, time.Duration, error) {
	if atomic.LoadUint32(&l.redisAlive) == 0 {
		return l.rescueLimiter.Allow(ctx, key)
	}

	res, err := l.limiter.Allow(ctx, key, redis_rate.Limit{
		Rate:   l.rate,
		Period: l.period,
		Burst:  l.burst,
	})
	if errors.Is(err, redis.Nil) {
		return false, 0, nil
	}
	if err != nil {
		l.startMonitor()
		return l.rescueLimiter.Allow(ctx, key)
	}

	return res.Allowed > 0, res.RetryAfter, nil
}

func (l *redisLimiter) Wait(ctx context.Context, key string) error {
	for {
		allowed, retryAfter, err := l.Allow(ctx, key)
		if err != nil {
			return err
		}
		if allowed {
			break
		}

		// 如果不允许，等待 RetryAfter 指定的时间再重试
		timer := time.NewTimer(retryAfter)
		select {
		case <-timer.C:
			// Timer expired, retry
			timer.Stop()
		case <-ctx.Done():
			// Context was cancelled, return error
			timer.Stop()
			return ctx.Err()
		}
	}
	return nil
}

func (l *redisLimiter) startMonitor() {
	l.rescueLock.Lock()
	defer l.rescueLock.Unlock()

	if l.monitorStarted {
		return
	}

	l.monitorStarted = true
	atomic.StoreUint32(&l.redisAlive, 0)

	go l.waitForRedis()
}

func (l *redisLimiter) waitForRedis() {
	ticker := time.NewTicker(pingInterval)
	defer func() {
		ticker.Stop()
		l.rescueLock.Lock()
		l.monitorStarted = false
		l.rescueLock.Unlock()
	}()

	for range ticker.C {
		if err := l.client.Ping(context.Background()).Err(); err == nil {
			atomic.StoreUint32(&l.redisAlive, 1)
			return
		}
	}
}
