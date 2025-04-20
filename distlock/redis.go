package distlock

import (
	"context"
	"sync"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

// RedisStorage 是基于 Redis 实现的 LockStorage
type RedisStorage struct {
	config Config
	mu     sync.Mutex
	rs     *redsync.Redsync // 这里持有一个 Redsync 实例
	mutex  *redsync.Mutex   // 互斥锁
}

// NewRedisStorage 创建一个新的 RedisStorage 实例
func NewRedisStorage(client goredislib.UniversalClient, config Config) *RedisStorage {
	rs := redsync.New(goredis.NewPool(client))
	mutex := rs.NewMutex(config.Key, redsync.WithExpiry(config.TTL))
	return &RedisStorage{
		config: config,
		rs:     rs,
		mutex:  mutex,
	}
}

// Lock 获取锁
func (r *RedisStorage) Lock(ctx context.Context) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	if err := r.mutex.LockContext(ctx); err != nil {
		return false, err
	}
	return true, nil
}

// Unlock 释放锁
func (r *RedisStorage) Unlock(ctx context.Context) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.mutex.UnlockContext(ctx)
}

// Renewal 锁续期
func (r *RedisStorage) Renewal(ctx context.Context) (bool, error) {
	r.mu.Lock()
	defer r.mu.Unlock()
	return r.mutex.ExtendContext(ctx)
}
