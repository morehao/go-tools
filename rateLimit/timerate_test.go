package rateLimit

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestTimeRateAllow(t *testing.T) {

	limiter := newTimeRateLimiter(time.Second, 1, time.Minute)

	ctx := context.Background()

	// Test Allow method
	for i := 0; i < 10; i++ {
		allowed, afterRetry, err := limiter.Allow(ctx, "test_key")
		assert.Nil(t, err)
		fmt.Println("Allow method - i:", i, "allowed:", allowed, "afterRetry:", afterRetry)
		time.Sleep(300 * time.Millisecond) // 每次请求之间间隔300毫秒
	}
}
