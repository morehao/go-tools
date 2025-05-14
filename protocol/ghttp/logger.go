package ghttp

import (
	"time"

	"github.com/morehao/go-tools/glog"
	"resty.dev/v3"
)

// LoggingMiddleware 返回一个日志中间件
func LoggingMiddleware(cfg *Client) func(c *resty.Client, resp *resty.Response) error {
	return func(c *resty.Client, resp *resty.Response) error {
		ctx := resp.Request.Context()
		begin := resp.Request.Time
		cost := glog.GetRequestCost(begin, time.Now())
		responseBody := resp.Result()
		fields := []any{
			glog.KeyProto, glog.ValueProtoHttp,
			glog.KeyService, cfg.Service,
			glog.KeyHost, cfg.Host,
			glog.KeyUri, resp.Request.URL,
			glog.KeyMethod, resp.Request.Method,
			glog.KeyHttpStatusCode, resp.StatusCode(),
			glog.KeyRequestBody, resp.Request.Body,
			glog.KeyRequestQuery, resp.Request.QueryParams,
			glog.KeyResponseBody, responseBody,
			glog.KeyCost, cost,
		}

		if resp.IsError() {
			// 记录错误日志
			fields = append(fields, glog.KeyErrorMsg, resp.Error())
			glog.Errorw(ctx, "HTTP request fail", fields...)
		} else {
			// 记录成功日志
			glog.Infow(ctx, "HTTP request success", fields...)
		}

		return nil
	}
}
