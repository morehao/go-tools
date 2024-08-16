package rateLimit

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

// TestTimeRateLimiter 测试基于内存的限流器
func TestTimeRateLimiter(t *testing.T) {
	ctx := context.Background()
	limiter, err := NewLimiter(
		WithMode(ModeTimeRate),
		WithPeriod(time.Second),
		WithRate(3),
		WithBurst(3),
		WithCleanupInterval(time.Minute),
	)
	assert.NoError(t, err)

	for i := 0; i < 10; i++ {
		allowed, err := limiter.Allow(ctx, "test_key")
		assert.Nil(t, err)
		fmt.Println("i:", i, "allowed:", allowed)
		// time.Sleep(time.Second)
	}
}

// TestRedisLimiter 测试基于 Redis 的限流器
func TestRedisLimiter(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	pong, err := client.Ping(context.Background()).Result()
	assert.Nil(t, err)
	fmt.Println("Redis connection successful:", pong)
	defer client.Close()

	limiter, err := NewLimiter(
		WithMode(ModeRedis),
		WithPeriod(time.Second),
		WithRate(3),
		WithBurst(3),
		WithRedisClient(client),
	)
	assert.NoError(t, err)

	ctx := context.Background()

	for i := 0; i < 10; i++ {
		allowed, err := limiter.Allow(ctx, "test_key")
		assert.Nil(t, err)
		fmt.Println("i:", i, "allowed:", allowed)
		// time.Sleep(time.Second)
	}
}
