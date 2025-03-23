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
	// clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
	// 	Addrs: []string{"127.0.0.1:6379"},
	// })
	rdbClient := redis.NewClient(&redis.Options{
		Addr: "127.0.0.1:6379",
	})
	config := Config{
		Key:         "test_lock",
		TTL:         time.Second * 5,
		AutoRenewal: true,
	}
	redisStore := NewRedisStorage(rdbClient, config)

	// 创建锁配置

	// 创建锁实例
	lock := NewDistLock(redisStore, &config)
	ctx := context.Background()
	// 获取锁
	firstLockRes, firstLockErr := lock.Lock(ctx)
	assert.Nil(t, firstLockErr)
	t.Log("first lock result: ", firstLockRes)
	// secondLockOk, secondLockErr := lock.Lock(ctx)
	// assert.Nil(t, secondLockErr)
	// t.Log("second lock result: ", secondLockOk)
	time.Sleep(time.Second * 10)

	unlockRes, unlockErr := lock.Unlock(ctx)
	assert.Nil(t, unlockErr)
	t.Log("unlockRes result: ", unlockRes)

}
