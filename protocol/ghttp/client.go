/*
 * @Author: morehao morehao@qq.com
 * @Date: 2025-05-14 10:07:23
 * @LastEditors: morehao morehao@qq.com
 * @LastEditTime: 2025-05-14 10:46:58
 * @FilePath: /go-tools/protocol/ghttp/client.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package ghttp

import (
	"context"
	"sync"
	"time"

	"github.com/morehao/go-tools/glog"
	"resty.dev/v3"
)

type Client struct {
	Module  string        `yaml:"service"`
	Host    string        `yaml:"host"`
	Timeout time.Duration `yaml:"timeout"`
	Retry   int           `yaml:"retry"`

	logger glog.Logger
	client *resty.Client
	once   sync.Once
}

// NewClient 创建一个新的 HTTP 客户端
func NewClient(cfg *Client) *Client {
	if cfg == nil {
		cfg = &Client{}
	}
	cfg.init()
	return cfg
}

// init 初始化客户端
func (c *Client) init() {
	c.once.Do(func() {
		// 初始化 HTTP 客户端
		client := resty.New()

		// 设置超时
		if c.Timeout > 0 {
			client.SetTimeout(c.Timeout)
		}

		// 设置重试
		if c.Retry > 0 {
			client.SetRetryCount(c.Retry)
		}

		// 设置基础配置
		client.SetHeader("Service", c.Module)
		if c.Host != "" {
			client.SetBaseURL(c.Host)
		}

		// 初始化 logger
		logCfg := glog.GetLoggerConfig()
		logCfg.Module = c.Module
		if logger, err := glog.GetLogger(logCfg); err != nil {
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
