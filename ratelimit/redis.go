package ratelimit

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis_rate/v10"
	"github.com/redis/go-redis/v9"
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

func (l *redisLimiter) Allow(ctx context.Context, key string) (bool, error) {
	if atomic.LoadUint32(&l.redisAlive) == 0 {
		return l.rescueLimiter.Allow(ctx, key), nil
	}

	res, err := l.limiter.Allow(ctx, key, redis_rate.Limit{
		Rate:   l.rate,
		Period: l.period,
		Burst:  l.burst,
	})
	if errors.Is(err, redis.Nil) {
		return false, nil
	}
	if err != nil {
		l.startMonitor()
		return l.rescueLimiter.Allow(ctx, key), nil
	}

	return res.Allowed > 0, nil
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
