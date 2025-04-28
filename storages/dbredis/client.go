package dbredis

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/morehao/go-tools/glog"
	"github.com/redis/go-redis/v9"
)

var (
	dbMap = map[string]*redis.Client{}
	lock  sync.RWMutex
)

type RedisConfig struct {
	Service      string        `yaml:"service"`       // 服务名
	Addr         string        `yaml:"addr"`          // redis地址
	Password     string        `yaml:"password"`      // 密码
	DB           int           `yaml:"db"`            // 数据库
	DialTimeout  time.Duration `yaml:"dial_timeout"`  // 连接超时
	ReadTimeout  time.Duration `yaml:"read_timeout"`  // 读取超时
	WriteTimeout time.Duration `yaml:"write_timeout"` // 写入超时
}

func InitRedis(cfg RedisConfig) (*redis.Client, error) {
	if cfg.Service == "" {
		return nil, fmt.Errorf("service name is empty")
	}
	opt := &redis.Options{
		Addr:             cfg.Addr,
		Password:         cfg.Password,
		DB:               cfg.DB,
		DisableIndentity: true,
	}
	if cfg.DialTimeout > 0 {
		opt.DialTimeout = cfg.DialTimeout
	}
	if cfg.ReadTimeout > 0 {
		opt.ReadTimeout = cfg.ReadTimeout
	}
	if cfg.WriteTimeout > 0 {
		opt.WriteTimeout = cfg.WriteTimeout
	}
	rdb := redis.NewClient(&redis.Options{
		Addr:             cfg.Addr,
		Password:         cfg.Password,
		DB:               cfg.DB,
		DisableIndentity: true,
	})
	service := cfg.Service
	if service == "" {
		service = "redis"
	}
	l, getLoggerErr := glog.GetModuleLogger("redis", glog.WithCallerSkip(6))
	if getLoggerErr != nil {
		return nil, getLoggerErr
	}
	logger := redisLogger{
		Service:  service,
		Addr:     cfg.Addr,
		Database: cfg.DB,
		Logger:   l,
	}
	rdb.AddHook(logger)
	// 发送PING命令，检查连接是否正常
	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	} else {
		fmt.Println("Redis connection successful:", pong)
	}
	lock.Lock()
	defer lock.Unlock()
	dbMap[cfg.Service] = rdb
	return rdb, nil
}

func InitMultiRedis(configs []RedisConfig) error {
	if len(configs) == 0 {
		return fmt.Errorf("redis config is empty")
	}
	for _, cfg := range configs {
		_, err := InitRedis(cfg)
		if err != nil {
			return err
		}
	}
	return nil
}

func GetClient(service string) *redis.Client {
	lock.RLock()
	defer lock.RUnlock()
	return dbMap[service]
}
