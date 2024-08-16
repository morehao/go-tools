package rateLimit

import (
	"context"
	"golang.org/x/time/rate"
	"sync"
	"time"
)

type timeRateLimiter struct {
	mu              sync.Mutex
	limiterMap      map[string]*rate.Limiter
	lastAccessedMap map[string]time.Time // 记录每个key的最后访问时间
	period          time.Duration        // 限制周期
	burst           int                  // 限制周期突发内允许的请求数
	cleanupInterval time.Duration        // 清理过期限流器的间隔
}

func (l *timeRateLimiter) Allow(ctx context.Context, key string) (bool, error) {
	l.mu.Lock()
	defer l.mu.Unlock()

	limiter, ok := l.limiterMap[key]
	if !ok {
		limiter = rate.NewLimiter(rate.Every(l.period), l.burst)
		l.limiterMap[key] = limiter
	}
	l.lastAccessedMap[key] = time.Now()

	return limiter.Allow(), nil
}

func (l *timeRateLimiter) Wait(ctx context.Context, key string) error {
	l.mu.Lock()
	defer l.mu.Unlock()
	limiter, exists := l.limiterMap[key]
	if !exists {
		limiter = rate.NewLimiter(rate.Every(l.period), l.burst)
		l.limiterMap[key] = limiter
	}
	l.lastAccessedMap[key] = time.Now()
	return limiter.Wait(ctx)
}

// 清理过期的限流器实例
func (l *timeRateLimiter) cleanupLoop() {
	for {
		time.Sleep(l.cleanupInterval)
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
