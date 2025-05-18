package ghttp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/morehao/golib/glog"
	"resty.dev/v3"
)

type SSEClientConfig struct {
	Module        string        `yaml:"service"`
	Host          string        `yaml:"host"`
	RetryWaitTime time.Duration `yaml:"retry_timeout"`
	Retry         int           `yaml:"retry"`
}

type SSEClient struct {
	Config SSEClientConfig
	logger glog.Logger
	es     *resty.EventSource
	once   sync.Once
}

// NewSSEClient 创建实例
func NewSSEClient(cfg *SSEClientConfig) *SSEClient {
	client := &SSEClient{
		Config: getDefaultSSEConfig(),
	}
	if cfg != nil {
		client = &SSEClient{
			Config: *cfg,
		}
	}
	client.init()
	return client
}

func (client *SSEClient) init() {
	client.once.Do(func() {
		es := resty.NewEventSource()
		if client.Config.RetryWaitTime > 0 {
			es.SetRetryWaitTime(client.Config.RetryWaitTime)
		}
		if client.Config.Retry > 0 {
			es.SetRetryCount(client.Config.Retry)
		}
		if client.Config.Module != "" {
			es.SetHeader("module", client.Config.Module)
		}

		logCfg := glog.GetLoggerConfig()
		logCfg.Module = client.Config.Module
		if logger, err := glog.GetLogger(logCfg, glog.WithCallerSkip(1)); err != nil {
			client.logger = glog.GetDefaultLogger()
		} else {
			client.logger = logger
		}

		client.es = es
	})
}

func (client *SSEClient) Es() *resty.EventSource {
	if client.es == nil {
		client.init()
	}
	return client.es
}

func (client *SSEClient) NewOpenHandler(ctx context.Context) resty.EventOpenFunc {
	return func(url string) {
		client.logger.Infow(ctx, "Http SSE Open",
			glog.KeyProto, glog.ValueProtoHttp,
			glog.KeyHost, client.Config.Host,
			glog.KeyUri, url,
		)
	}
}

// NewErrorHandler 构造带 context 的 OnError 日志函数
func (client *SSEClient) NewErrorHandler(ctx context.Context) resty.EventErrorFunc {
	return func(err error) {
		client.logger.Errorw(ctx, "Http SSE Error",
			glog.KeyProto, glog.ValueProtoHttp,
			glog.KeyHost, client.Config.Host,
			"error", err,
		)
	}
}

func (client *SSEClient) NewMessageHandler(ctx context.Context) resty.EventMessageFunc {
	return func(e any) {
		ev, ok := e.(*resty.Event)
		if !ok {
			client.logger.Errorw(ctx, "Invalid SSE message type", "type", fmt.Sprintf("%T", e))
			return
		}
		fmt.Println("ID:", ev.ID, "Name:", ev.Name, "Data:", ev.Data)
	}
}

func getDefaultSSEConfig() SSEClientConfig {
	return SSEClientConfig{
		Module: "httpSSE",
		Retry:  3,
	}
}
