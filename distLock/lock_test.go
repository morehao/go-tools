package distLock

import (
	"context"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
)

func TestLock(t *testing.T) {
	// 初始化Redis存储
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
	})
	config := Config{
		Key:   "test",
		TTL:   time.Second * 5,
		Owner: GenerateOwner(),
	}
	redisStore := NewRedisStorage([]*redis.Client{redisClient}, config)

	// 创建锁配置

	// 创建锁实例
	lock := NewDistLock(redisStore, &config)
	ctx := context.Background()
	// 获取锁
	firstLockRes, firstLockErr := lock.Lock(ctx)
	assert.Nil(t, firstLockErr)
	t.Log("first lock result: ", firstLockRes)
	secondLockOk, secondLockErr := lock.Lock(ctx)
	assert.Nil(t, secondLockErr)
	t.Log("second lock result: ", secondLockOk)

	unlockRes, unlockErr := lock.Unlock(context.Background())
	assert.Nil(t, unlockErr)
	t.Log("unlockRes result: ", unlockRes)

}
