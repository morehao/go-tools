/*
 * @Author: morehao morehao@qq.com
 * @Date: 2025-05-14 10:46:54
 * @LastEditors: morehao morehao@qq.com
 * @LastEditTime: 2025-05-14 12:13:44
 * @FilePath: /golib/protocol/ghttp/client_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package ghttp

import (
	"context"
	"testing"
	"time"

	"github.com/morehao/golib/glog"
	"github.com/stretchr/testify/assert"
)

func TestRWithResult(t *testing.T) {
	cfg := &ClientConfig{
		Module:  "httpbin",
		Host:    "http://httpbin.org",
		Timeout: 5 * time.Second,
		Retry:   3,
	}
	client := NewClient(cfg)
	ctx := context.Background()
	type Result struct {
		Args struct {
			Name string `json:"name"`
		} `json:"args"`
	}
	var result Result
	_, err := client.RWithResult(ctx, &result).SetQueryParam("name", "张三").Get("/get")

	assert.Nil(t, err)
	t.Log(glog.ToJsonString(result))
}

func TestRWithResultWithoutNew(t *testing.T) {
	cfg := ClientConfig{
		Module:  "httpbin",
		Host:    "http://httpbin.org",
		Timeout: 5 * time.Second,
		Retry:   3,
	}
	client := &Client{
		Config: cfg,
	}

	ctx := context.Background()
	type Result struct {
		Args struct {
			Name string `json:"name"`
		} `json:"args"`
	}
	var result Result
	_, err := client.RWithResult(ctx, &result).SetQueryParam("name", "张三").Get("/get")

	assert.Nil(t, err)
	t.Log(glog.ToJsonString(result))
}

func TestMultiClient(t *testing.T) {
	client1 := &Client{
		Config: ClientConfig{
			Module:  "httpbin1",
			Host:    "http://httpbin.org",
			Timeout: 5 * time.Second,
			Retry:   3,
		},
	}
	client2 := &Client{
		Config: ClientConfig{
			Module:  "httpbin2",
			Host:    "http://httpbin.org",
			Timeout: 5 * time.Second,
			Retry:   3,
		},
	}
	ctx := context.Background()
	type Result struct {
		Args struct {
			Name string `json:"name"`
		} `json:"args"`
	}
	var result1 Result
	_, err := client1.RWithResult(ctx, &result1).SetQueryParam("name", "张三").Get("/get")

	assert.Nil(t, err)
	t.Log(glog.ToJsonString(result1))
	var result2 Result
	_, err2 := client2.RWithResult(ctx, &result2).SetQueryParam("name", "李四").Get("/get")

	assert.Nil(t, err2)
	t.Log(glog.ToJsonString(result2))
}
