package dbredis

import (
	"context"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/morehao/go-tools/glog"
	"github.com/redis/go-redis/v9"
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
	l, getLoggerErr := glog.GetModuleLogger("redis", glog.WithCallerSkip(1))
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
	return rdb, nil
}

type redisLogger struct {
	Service  string
	Addr     string
	Database int
	Logger   glog.Logger
}

// DialHook 当创建网络连接时调用的hook
func (l redisLogger) DialHook(next redis.DialHook) redis.DialHook {
	return func(ctx context.Context, network, addr string) (net.Conn, error) {
		return next(ctx, network, addr)
	}
}

// ProcessHook 执行命令时调用的hook
func (l redisLogger) ProcessHook(next redis.ProcessHook) redis.ProcessHook {
	return func(ctx context.Context, cmd redis.Cmder) error {

		begin := time.Now()
		fields := l.commonFields(ctx)
		fields = append(fields,
			glog.KeyCmd, cmd.FullName(),
		)
		var ralCode int
		if err := cmd.Err(); err != nil {
			msg := err.Error()
			ralCode = -1
			end := time.Now()
			cost := glog.GetRequestCost(begin, end)
			fields = append(fields,
				glog.KeyCmdContent, cmd.String(),
				glog.KeyRalCode, ralCode,
				glog.KeyCost, cost,
			)
			l.Logger.Errorw(ctx, msg, fields...)
			return err
		}

		hook := next(ctx, cmd)

		end := time.Now()
		cost := glog.GetRequestCost(begin, end)
		fields = append(fields,
			glog.KeyCmdContent, cmd.String(),
			glog.KeyRalCode, ralCode,
			glog.KeyCost, cost,
		)

		l.Logger.Debugw(ctx, "redis execute success", fields...)
		return hook
	}
}

// ProcessPipelineHook 执行管道命令时调用的hook
func (l redisLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		begin := time.Now() // 记录开始时间
		err := next(ctx, cmds)
		end := time.Now() // 记录结束时间
		cost := glog.GetRequestCost(begin, end)

		// 准备日志字段
		fields := l.commonFields(ctx)
		fields = append(fields,
			glog.KeyCmdContent, l.cmdsToString(cmds),
			glog.KeyCost, cost,
		)

		// 根据执行结果记录日志
		if err != nil {
			fields = append(fields, glog.KeyRalCode, -1)
			l.Logger.Errorw(ctx, fmt.Sprintf("redis pipeline execute failed, err: %v", err), fields...)
		} else {
			fields = append(fields, glog.KeyRalCode, 0)
			l.Logger.Debugw(ctx, "redis pipeline execute success", fields...)
		}
		return err
	}
}

// cmdsToString 将管道命令转换为字符串表示，用于日志记录
func (l redisLogger) cmdsToString(cmds []redis.Cmder) string {
	var cmdStrs []string
	for _, cmd := range cmds {
		cmdStrs = append(cmdStrs, cmd.String())
	}
	return fmt.Sprintf("[%s]", strings.Join(cmdStrs, ", "))
}
func (l redisLogger) commonFields(ctx context.Context) []interface{} {
	fields := []interface{}{
		glog.KeyAddr, l.Addr,
		glog.KeyDatabase, l.Database,
	}
	return fields
}
