package dbClient

import (
	"context"
	"github.com/morehao/go-tools/glog"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"net"
	"time"
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

func InitRedis(cfg RedisConfig) *redis.Client {
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
	logger := redisLogger{
		Service:  service,
		Addr:     cfg.Addr,
		Database: cfg.DB,
		Logger:   glog.GetLogger(glog.WithZapOptions(zap.AddCallerSkip(3))),
	}
	rdb.AddHook(logger)
	return rdb
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
			"cmd", cmd.FullName(),
		)
		var ralCode int
		if err := cmd.Err(); err != nil {
			msg := err.Error()
			ralCode = -1
			end := time.Now()
			cost := glog.GetRequestCost(begin, end)
			fields = append(fields,
				"cmdContent", cmd.String(),
				"ralCode", ralCode,
				"requestStartTime", glog.FormatRequestTime(begin),
				"requestEndTime", glog.FormatRequestTime(end),
				"cost", cost,
			)
			l.Logger.Errorw(ctx, msg, fields...)
			return err
		}

		hook := next(ctx, cmd)

		end := time.Now()
		cost := glog.GetRequestCost(begin, end)
		fields = append(fields,
			"cmdContent", cmd.String(),
			"ralCode", ralCode,
			"requestStartTime", glog.FormatRequestTime(begin),
			"requestEndTime", glog.FormatRequestTime(end),
			"cost", cost,
		)

		l.Logger.Infow(ctx, "redis execute success", fields...)
		return hook
	}
}

// ProcessPipelineHook 执行管道命令时调用的hook
func (l redisLogger) ProcessPipelineHook(next redis.ProcessPipelineHook) redis.ProcessPipelineHook {
	return func(ctx context.Context, cmds []redis.Cmder) error {
		return next(ctx, cmds)
	}
}

func (l *redisLogger) commonFields(ctx context.Context) []interface{} {
	fields := []interface{}{
		glog.KeyProto, "redis",
		"service", l.Service,
		"addr", l.Addr,
		"database", l.Database,
	}
	return fields
}
