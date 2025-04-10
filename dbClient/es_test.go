package dbClient

import (
	"context"
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
	query := &types.Query{
		Match: map[string]types.MatchQuery{
			"firstname": {
				Query: "Amber",
			},
		},
	}

	res, searchErr := typedClient.Search().
		Index("accounts").
		Query(query).Do(ctx)
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
	res, searchErr := simpleClient.Get("accounts", "1")
	assert.Nil(t, searchErr)
	t.Log(gutils.ToJsonString(res))
}
