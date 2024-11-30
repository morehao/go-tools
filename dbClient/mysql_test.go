package dbClient

import (
	"context"
	"testing"

	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitMysql(t *testing.T) {
	defer glog.Close()
	logCfg := &glog.LoggerConfig{
		Service:   "test",
		Level:     glog.InfoLevel,
		Dir:       "./log",
		Stdout:    true,
		ExtraKeys: []string{"requestId"},
	}
	opt := glog.WithZapOptions(zap.AddCallerSkip(3))
	initLogErr := glog.NewLogger(logCfg, opt)
	assert.Nil(t, initLogErr)
	cfg := MysqlConfig{
		Service:  "test",
		Addr:     "127.0.0.1:3306",
		Database: "demo",
		User:     "root",
		Password: "123456",
	}
	mysqlClient, initDbErr := InitMysql(cfg)
	assert.Nil(t, initDbErr)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", "12312312312312")
	var res []interface{}
	findErr := mysqlClient.WithContext(ctx).Table("user").Find(&res).Error
	assert.Nil(t, findErr)
	t.Log(gutils.ToJsonString(res))
}

func TestInitMysqlWithoutInitLog(t *testing.T) {
	defer glog.Close()
	cfg := MysqlConfig{
		Service:  "test",
		Addr:     "127.0.0.1:3306",
		Database: "demo",
		User:     "root",
		Password: "123456",
	}
	mysqlClient, initDbErr := InitMysql(cfg)
	assert.Nil(t, initDbErr)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", "12312312312312")
	var res []interface{}
	findErr := mysqlClient.WithContext(ctx).Table("user").Find(&res).Error
	assert.Nil(t, findErr)
	t.Log(gutils.ToJsonString(res))
}
