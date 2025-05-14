package ghttp

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/morehao/go-tools/glog"
	"resty.dev/v3"
)

type SSEInstConfig struct {
	Module        string        `yaml:"service"`
	Host          string        `yaml:"host"`
	RetryWaitTime time.Duration `yaml:"retry_timeout"`
	Retry         int           `yaml:"retry"`
}

type SSEInst struct {
	Config SSEInstConfig
	logger glog.Logger
	es     *resty.EventSource
	once   sync.Once
}

// NewSSEInst 创建实例
func NewSSEInst(cfg *SSEInstConfig) *SSEInst {
	client := &SSEInst{
		Config: getDefaultSSEConfig(),
	}
	if cfg != nil {
		client = &SSEInst{
			Config: *cfg,
		}
	}
	client.init()
	return client
}

func (inst *SSEInst) init() {
	inst.once.Do(func() {
		es := resty.NewEventSource()
		if inst.Config.RetryWaitTime > 0 {
			es.SetRetryWaitTime(inst.Config.RetryWaitTime)
		}
		if inst.Config.Retry > 0 {
			es.SetRetryCount(inst.Config.Retry)
		}
		if inst.Config.Module != "" {
			es.SetHeader("module", inst.Config.Module)
		}

		logCfg := glog.GetLoggerConfig()
		logCfg.Module = inst.Config.Module
		if logger, err := glog.GetLogger(logCfg); err != nil {
			inst.logger = glog.GetDefaultLogger()
		} else {
			inst.logger = logger
		}

		inst.es = es
	})
}

func (inst *SSEInst) Es() *resty.EventSource {
	if inst.es == nil {
		inst.init()
	}
	return inst.es
}

func (inst *SSEInst) NewOpenHandler(ctx context.Context) resty.EventOpenFunc {
	return func(url string) {
		inst.logger.Infow(ctx, "Http SSE Open",
			glog.KeyProto, glog.ValueProtoHttp,
			glog.KeyHost, inst.Config.Host,
			glog.KeyUri, url,
		)
	}
}

// NewErrorHandler 构造带 context 的 OnError 日志函数
func (inst *SSEInst) NewErrorHandler(ctx context.Context) resty.EventErrorFunc {
	return func(err error) {
		inst.logger.Errorw(ctx, "Http SSE Error",
			glog.KeyProto, glog.ValueProtoHttp,
			glog.KeyHost, inst.Config.Host,
			"error", err,
		)
	}
}

func (inst *SSEInst) NewMessageHandler(ctx context.Context) resty.EventMessageFunc {
	return func(e any) {
		ev, ok := e.(*resty.Event)
		if !ok {
			inst.logger.Errorw(ctx, "Invalid SSE message type", "type", fmt.Sprintf("%T", e))
			return
		}
		fmt.Println("ID:", ev.ID, "Name:", ev.Name, "Data:", ev.Data)
	}
}

func getDefaultSSEConfig() SSEInstConfig {
	return SSEInstConfig{
		Module: "httpSSE",
		Retry:  3,
	}
}
