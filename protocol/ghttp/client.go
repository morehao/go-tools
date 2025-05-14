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
func NewClient(cfg *Client) (*Client, error) {
	if cfg == nil {
		cfg = &Client{}
	}
	if err := cfg.init(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// init 初始化客户端
func (c *Client) init() error {
	var err error
	c.once.Do(func() {
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

		// 添加日志中间件
		client.AddResponseMiddleware(LoggingMiddleware(c))

		// 设置日志
		logCfg := glog.GetLoggerConfig()
		logCfg.Module = c.Module
		logger, getLoggerErr := glog.GetLogger(logCfg)
		if getLoggerErr != nil {
			err = getLoggerErr
			return
		}
		c.logger = logger
		c.client = client
	})
	return err
}

// R 创建一个新的请求，支持 context
func (c *Client) R(ctx context.Context) (*resty.Request, error) {
	if c.client == nil {
		if err := c.init(); err != nil {
			return nil, err
		}
	}
	return c.client.R().SetContext(ctx), nil
}

func (c *Client) RWithResult(ctx context.Context, result any) (*resty.Request, error) {
	if c.client == nil {
		if err := c.init(); err != nil {
			return nil, err
		}
	}
	return c.client.R().
		SetContext(ctx).
		SetResult(result), nil
}
