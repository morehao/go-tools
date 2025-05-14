package ghttp

import (
	"context"
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

func (inst *SSEInst) Get() {

}

func (inst *SSEInst) OnOpen(ctx context.Context) resty.EventOpenFunc {
	return func(url string) {
		fields := []any{
			glog.KeyProto, glog.ValueProtoHttp,
			glog.KeyHost, inst.Config.Host,
			glog.KeyUri, url,
			glog.KeyMethod, "",
			glog.KeyHttpStatusCode, resp.StatusCode(),
			glog.KeyRequestBody, resp.Request.Body,
			glog.KeyRequestQuery, resp.Request.QueryParams.Encode(),
			glog.KeyResponseBody, responseBody,
			glog.KeyCost, cost,
		}
		inst.logger.Infow(ctx)
	}
}

func getDefaultSSEConfig() SSEInstConfig {
	return SSEInstConfig{
		Module: "httpSSE",
		Retry:  3,
	}
}
