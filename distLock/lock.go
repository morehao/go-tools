package distLock

import (
	"context"
	"fmt"
	"sync"
	"time"
)

// Lock 锁接口（支持不同存储引擎扩展）
// 暂不考虑可重入性，个人理解，如果存在重复获取锁，通过代码逻辑调整即可实现
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
}

type DistLock struct {
	store  Lock
	config *Config
	// count    int          // 可重入计数器
	mu       sync.RWMutex // 使用 RWMutex 提高并发读性能
	stopOnce sync.Once    // 确保 stopChan 只关闭一次
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
	// if dl.count > 0 {
	// 	dl.count++
	// 	return true, nil
	// }

	// 获取新锁
	ok, err := dl.store.Lock(ctx)
	if err != nil {
		return false, err
	}
	if !ok {
		return false, fmt.Errorf("lock acquisition failed")
	}

	// dl.count = 1

	// 启动自动续期
	if dl.config.AutoRenewal {
		go dl.autoRenewal(ctx)
	}

	return true, nil
}

func (dl *DistLock) Unlock(ctx context.Context) (bool, error) {
	dl.mu.Lock()
	defer dl.mu.Unlock()

	// if dl.count == 0 {
	// 	return false, fmt.Errorf("lock not held")
	// }
	//
	// dl.count--
	// if dl.count > 0 {
	// 	return true, nil
	// }

	// 停止续期
	if dl.config.AutoRenewal {
		dl.stopOnce.Do(func() {
			close(dl.stopChan)
		})
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
