package glog

import (
	"context"
	"reflect"
	"strconv"
	"time"
)

func GetTraceId(ctx context.Context) string {
	if ctx == nil {
		return genTraceID()
	}
	// 如果上下文有requestId，则直接返回
	if traceIdVal := ctx.Value(ContextKeyTraceId); traceIdVal != nil {
		if traceId := traceIdVal.(string); traceId != "" {
			return traceId
		}
		return genTraceID()
	}

	traceId := genTraceID()
	return traceId
}

// genTraceID 生成请求ID,TODO: 生成规则待定
func genTraceID() (traceId string) {
	// 生成纳秒时间戳
	nanosecond := uint64(time.Now().UnixNano())
	// nanosecond&0x7FFFFFFF 使用位运算与操作，将 nanosecond 的二进制表示的最高位（最高位是符号位）清零，将其转换为正整数。
	// |0x80000000 使用位运算或操作，将二进制表示的最高位设置为 1，以确保结果是一个正整数。这样做的目的是为了确保结果是正数，而不是负数。
	traceId = strconv.FormatUint(nanosecond&0x7FFFFFFF|0x80000000, 10)
	return traceId
}

func GetIp(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if ip, ok := ctx.Value(ContextKeyIp).(string); ok {
		return ip
	}
	return ""
}

func GetUri(ctx context.Context) string {
	if ctx == nil {
		return ""
	}
	if uri, ok := ctx.Value(ContextKeyUri).(string); ok {
		return uri
	}
	return ""
}

func GetFormatRequestTime(time time.Time) string {
	return time.Format("2006-01-02 15:04:05.999999")
}

func GetRequestCost(start, end time.Time) float64 {
	return float64(end.Sub(start).Nanoseconds()/1e4) / 100.0
}

func nilCtx(ctx context.Context) bool {
	return ctx == nil
}

func IsNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
