package distLock

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Lock 锁接口（支持不同存储引擎扩展）
type Lock interface {
	Lock(ctx context.Context) (bool, error)
	Unlock(ctx context.Context) (bool, error)
	Renewal(ctx context.Context) (bool, error)
}

// Config 锁配置
type Config struct {
	AutoRenewal bool          // 是否自动续期
	TTL         time.Duration // 锁的超时时间
	Key         string
	Owner       string // 锁的拥有者标识，防止其他客户端误释放锁
}

type DistLock struct {
	store    Lock
	config   *Config
	count    int        // 可重入计数器
	mu       sync.Mutex // 本地锁
	stopChan chan struct{}
}

// NewDistLock 创建新锁实例
func NewDistLock(store Lock, config *Config) *DistLock {

	return &DistLock{
		store:    store,
		config:   config,
		stopChan: make(chan struct{}),
	}
}

func (dl *DistLock) Lock(ctx context.Context) (bool, error) {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	// 可重入检查
	if dl.count > 0 {
		dl.count++
		return true, nil
	}

	// 获取新锁
	ok, err := dl.store.Lock(ctx)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("lock acquisition failed")
	}

	dl.count = 1

	// 启动自动续期
	if dl.config.AutoRenewal {
		go dl.autoRenewal(ctx)
	}

	return true, nil
}

func (dl *DistLock) Unlock(ctx context.Context) (bool, error) {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	if dl.count == 0 {
		return false, fmt.Errorf("lock not held")
	}

	dl.count--
	if dl.count > 0 {
		return true, nil
	}

	// 停止续期
	if dl.config.AutoRenewal {
		close(dl.stopChan)
	}

	// 释放锁
	return dl.store.Unlock(ctx)
}

// 自动续期循环
func (dl *DistLock) autoRenewal(ctx context.Context) {
	renewalInterval := dl.config.TTL / 2
	ticker := time.NewTicker(renewalInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			ok, err := dl.store.Renewal(ctx)

			if err != nil || !ok {
				return
			}
		case <-dl.stopChan:
			return
		}
	}
}
