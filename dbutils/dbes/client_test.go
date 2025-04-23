package dbes

import (
	"context"
	"strings"
	"testing"

	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInitTypedES(t *testing.T) {
	defer glog.Close()
	logCfg := &glog.LoggerConfig{
		Service:   "ES",
		Level:     glog.InfoLevel,
		Dir:       "./log",
		Stdout:    true,
		ExtraKeys: []string{"requestId"},
	}
	opt := glog.WithZapOptions(zap.AddCallerSkip(2))
	initLogErr := glog.NewLogger(logCfg, opt)
	assert.Nil(t, initLogErr)
	cfg := ESConfig{
		Service: "es",
		Addr:    "http://localhost:9200",
	}
	_, typedClient, initErr := InitES(cfg)
	assert.Nil(t, initErr)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", "12312312312312")

	res, searchErr := typedClient.Search().
		Index("accounts").
		Query(&types.Query{
			MatchAll: types.NewMatchAllQuery(),
		}).Do(ctx)
	assert.Nil(t, searchErr)
	t.Log(gutils.ToJsonString(res))
}

func TestInitSimpleES(t *testing.T) {
	defer glog.Close()
	logCfg := &glog.LoggerConfig{
		Service:   "ES",
		Level:     glog.InfoLevel,
		Dir:       "./log",
		Stdout:    true,
		ExtraKeys: []string{"requestId"},
	}
	opt := glog.WithZapOptions(zap.AddCallerSkip(2))
	initLogErr := glog.NewLogger(logCfg, opt)
	assert.Nil(t, initLogErr)
	cfg := ESConfig{
		Service: "es",
		Addr:    "http://localhost:9200",
	}
	simpleClient, _, initErr := InitES(cfg)
	assert.Nil(t, initErr)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", "12312312312312")
	res, searchErr := simpleClient.Search(
		simpleClient.Search.WithContext(ctx),
		simpleClient.Search.WithIndex("accounts"),
		simpleClient.Search.WithBody(strings.NewReader(`{"query":{"match_all":{}}}`)),
	)
	assert.Nil(t, searchErr)
	t.Log(gutils.ToJsonString(res))
}
