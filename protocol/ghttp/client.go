package ghttp

import (
	"context"
	"sync"
	"time"

	"github.com/morehao/go-tools/glog"
	"resty.dev/v3"
)

type ClientConfig struct {
	Module  string        `yaml:"module"`
	Host    string        `yaml:"host"`
	Timeout time.Duration `yaml:"timeout"`
	Retry   int           `yaml:"retry"`
}

type Client struct {
	Config ClientConfig
	logger glog.Logger
	client *resty.Client
	once   sync.Once
}

// NewClient 创建一个新的 HTTP 客户端
func NewClient(cfg *ClientConfig) *Client {
	client := &Client{
		Config: getDefaultConfig(),
	}

	if cfg != nil {
		client = &Client{
			Config: *cfg,
		}
	}
	client.init()
	return client
}

// init 初始化客户端
func (c *Client) init() {
	c.once.Do(func() {
		// 初始化 HTTP 客户端
		client := resty.New()

		// 设置超时
		if c.Config.Timeout > 0 {
			client.SetTimeout(c.Config.Timeout)
		}

		// 设置重试
		if c.Config.Retry > 0 {
			client.SetRetryCount(c.Config.Retry)
		}

		// 设置基础配置
		if c.Config.Module != "" {
			client.SetHeader("module", c.Config.Module)
		}
		if c.Config.Host != "" {
			client.SetBaseURL(c.Config.Host)
		}

		// 初始化 logger
		logCfg := glog.GetLoggerConfig()
		logCfg.Module = c.Config.Module
		if logger, err := glog.GetLogger(logCfg, glog.WithCallerSkip(1)); err != nil {
			c.logger = glog.GetDefaultLogger()
		} else {
			c.logger = logger
		}

		// 添加日志中间件
		client.AddResponseMiddleware(LoggingMiddleware(c))

		c.client = client
	})
}

// R 创建一个新的请求，支持 context
func (c *Client) R(ctx context.Context) *resty.Request {
	if c.client == nil {
		c.init()
	}
	return c.client.R().SetContext(ctx)
}

func (c *Client) RWithResult(ctx context.Context, result any) *resty.Request {
	if c.client == nil {
		c.init()
	}
	return c.client.R().
		SetContext(ctx).
		SetResult(result)
}

func getDefaultConfig() ClientConfig {
	return ClientConfig{
		Module: "httpClient",
		Retry:  3,
	}
}

// LoggingMiddleware 返回一个日志中间件
func LoggingMiddleware(client *Client) func(restyClient *resty.Client, resp *resty.Response) error {
	return func(c *resty.Client, resp *resty.Response) error {
		ctx := resp.Request.Context()
		begin := resp.Request.Time
		cost := glog.GetRequestCost(begin, time.Now())
		responseBody := resp.Result()
		fields := []any{
			glog.KeyProto, glog.ValueProtoHttp,
			glog.KeyHost, client.Config.Host,
			glog.KeyUri, resp.Request.URL,
			glog.KeyMethod, resp.Request.Method,
			glog.KeyHttpStatusCode, resp.StatusCode(),
			glog.KeyRequestBody, resp.Request.Body,
			glog.KeyRequestQuery, resp.Request.QueryParams.Encode(),
			glog.KeyResponseBody, responseBody,
			glog.KeyCost, cost,
		}

		if resp.IsError() {
			// 记录错误日志
			fields = append(fields, glog.KeyErrorMsg, resp.Error())
			client.logger.Errorw(ctx, "HTTP request fail", fields...)
		} else {
			// 记录成功日志
			client.logger.Infow(ctx, "HTTP request success", fields...)
		}

		return nil
	}
}
