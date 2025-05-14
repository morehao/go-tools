/*
 * @Author: morehao morehao@qq.com
 * @Date: 2025-05-14 10:46:54
 * @LastEditors: morehao morehao@qq.com
 * @LastEditTime: 2025-05-14 12:13:44
 * @FilePath: /go-tools/protocol/ghttp/client_test.go
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
package ghttp

import (
	"context"
	"testing"
	"time"

	"github.com/morehao/go-tools/glog"
	"github.com/stretchr/testify/assert"
)

func TestRWithResult(t *testing.T) {
	cfg := &Client{
		Module:  "httpbin",
		Host:    "http://httpbin.org",
		Timeout: 5 * time.Second,
		Retry:   3,
	}
	client, newErr := NewClient(cfg)
	assert.Nil(t, newErr)
	ctx := context.Background()
	type Result struct {
		Args struct {
			Name string `json:"name"`
		} `json:"args"`
	}
	var result Result
	request, newRequestErr := client.RWithResult(ctx, &result)
	assert.Nil(t, newRequestErr)
	_, err := request.SetQueryParam("name", "张三").Get("/get")

	assert.Nil(t, err)
	t.Log(glog.ToJsonString(result))
}

func TestWithoutNew(t *testing.T) {
	client := &Client{
		Module:  "httpbin",
		Host:    "http://httpbin.org",
		Timeout: 5 * time.Second,
		Retry:   3,
	}

	ctx := context.Background()
	type Result struct {
		Args struct {
			Name string `json:"name"`
		} `json:"args"`
	}
	var result Result
	request, newRequestErr := client.RWithResult(ctx, &result)
	assert.Nil(t, newRequestErr)
	_, err := request.SetQueryParam("name", "张三").Get("/get")

	assert.Nil(t, err)
	t.Log(glog.ToJsonString(result))
}
