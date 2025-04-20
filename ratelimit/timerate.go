package ratelimit

import (
	"context"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

type timeRateLimiter struct {
	mu              sync.Mutex
	limiterMap      map[string]*rate.Limiter
	lastAccessedMap map[string]time.Time // 记录每个key的最后访问时间
	period          time.Duration        // 限制周期
	burst           int                  // 限制周期突发内允许的请求数
	cleanupInterval time.Duration        // 清理过期限流器的间隔
}

// newTimeRateLimiter 创建一个新的 timeRateLimiter
func newTimeRateLimiter(period time.Duration, burst int, cleanupInterval time.Duration) *timeRateLimiter {
	limiter := &timeRateLimiter{
		limiterMap:      make(map[string]*rate.Limiter),
		lastAccessedMap: make(map[string]time.Time),
		period:          period,
		burst:           burst,
		cleanupInterval: cleanupInterval,
	}

	go limiter.cleanupLoop()

	return limiter
}

func (l *timeRateLimiter) Allow(ctx context.Context, key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, ok := l.limiterMap[key]
	if !ok {
		limiter = rate.NewLimiter(rate.Every(l.period), l.burst)
		l.limiterMap[key] = limiter
	}
	l.lastAccessedMap[key] = time.Now()

	return limiter.Allow()
}

// 清理过期的限流器实例
func (l *timeRateLimiter) cleanupLoop() {
	for range time.Tick(l.cleanupInterval) {
		l.cleanupExpiredLimiters()
	}
}

func (l *timeRateLimiter) cleanupExpiredLimiters() {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	for key, lastAccessed := range l.lastAccessedMap {
		if now.Sub(lastAccessed) > l.cleanupInterval {
			delete(l.limiterMap, key)
			delete(l.lastAccessedMap, key)
		}
	}
}
