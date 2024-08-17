package rateLimit

import (
	"github.com/redis/go-redis/v9"
	"time"
)

type Config struct {
	RedisClient     *redis.Client // redis 客户端
	Period          time.Duration // 限流周期
	CleanupInterval time.Duration // 清理过期限流器的间隔，只在ModeRateLimit模型下使用
	Rate            int           // 每个限流周期允许的最大请求数
	Burst           int           // 令牌桶的最大容量
}

type Option func(*Config)

func WithRedisClient(client *redis.Client) Option {
	return func(cfg *Config) {
		cfg.RedisClient = client
	}
}

// WithPeriod 设置限流周期
func WithPeriod(period time.Duration) Option {
	return func(cfg *Config) {
		cfg.Period = period
	}
}

// WithCleanupInterval 设置清理过期限流器的间隔
func WithCleanupInterval(cleanupInterval time.Duration) Option {
	return func(cfg *Config) {
		cfg.CleanupInterval = cleanupInterval
	}
}

// WithRate 设置限流周期内允许的最大请求数
func WithRate(rate int) Option {
	return func(cfg *Config) {
		cfg.Rate = rate
	}
}

// WithBurst 设置令牌桶的最大容量
func WithBurst(burst int) Option {
	return func(cfg *Config) {
		cfg.Burst = burst
	}
}
