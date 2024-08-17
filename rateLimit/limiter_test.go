package rateLimit

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestAllow(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	// pong, err := client.Ping(context.Background()).Result()
	// assert.Nil(t, err)
	// fmt.Println("Redis connection successful:", pong)
	defer client.Close()

	limiter, err := NewLimiter(
		WithPeriod(time.Second),
		WithRate(1),
		WithBurst(1),
		WithRedisClient(client),
	)
	assert.NoError(t, err)

	ctx := context.Background()

	// Test Allow method
	for i := 0; i < 10; i++ {
		allowed, afterRetry, err := limiter.Allow(ctx, "test_key")
		assert.Nil(t, err)
		fmt.Println("Allow method - i:", i, "allowed:", allowed, "afterRetry:", afterRetry)
		time.Sleep(300 * time.Millisecond) // 每次请求之间间隔300毫秒
	}
}

func TestWait(t *testing.T) {

	client := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
		DB:   0,
	})
	pong, err := client.Ping(context.Background()).Result()
	assert.Nil(t, err)
	fmt.Println("Redis connection successful:", pong)
	defer client.Close()

	limiter, err := NewLimiter(
		WithPeriod(time.Second),
		WithRate(1),
		WithBurst(1),
		WithRedisClient(client),
	)
	assert.NoError(t, err)

	ctx := context.Background()

	// Test Wait method
	for i := 0; i < 10; i++ {
		start := time.Now()
		err := limiter.Wait(ctx, "test_key")
		assert.Nil(t, err)
		duration := time.Since(start)
		fmt.Println("Wait method - i:", i, "waited for:", duration)
	}
}
