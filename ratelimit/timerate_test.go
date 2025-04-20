package ratelimit

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestTimeRateAllow(t *testing.T) {

	limiter := newTimeRateLimiter(time.Second, 1, time.Minute)

	ctx := context.Background()

	// Test Allow method
	for i := 0; i < 10; i++ {
		allowed := limiter.Allow(ctx, "test_key")
		fmt.Println("Allow method - i:", i, "allowed:", allowed)
		time.Sleep(300 * time.Millisecond) // 每次请求之间间隔300毫秒
	}
}
