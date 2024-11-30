package rateLimit

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
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
		allowed, err := limiter.Allow(ctx, "test_key")
		assert.Nil(t, err)
		fmt.Println("Allow method - i:", i, "allowed:", allowed)
		time.Sleep(300 * time.Millisecond) // 每次请求之间间隔300毫秒
	}
}
