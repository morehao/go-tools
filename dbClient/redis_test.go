package dbClient

import (
	"context"
	"github.com/morehao/go-tools/glog"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitRedis(t *testing.T) {
	defer glog.Close()
	err := glog.InitZapLogger(&glog.LoggerConfig{
		Service:   "test",
		Level:     glog.DebugLevel,
		Dir:       "./log",
		Stdout:    true,
		ExtraKeys: []string{"requestId"},
	})
	assert.Nil(t, err)
	cfg := RedisConfig{
		Service:  "test",
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
	redisClient := InitRedis(cfg)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", "12312312312312")
	if err := redisClient.Get(ctx, "test123").Err(); err != nil {
		assert.Nil(t, err)
	}
}
