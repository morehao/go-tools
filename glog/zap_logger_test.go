package glog

import (
	"context"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"regexp"
	"testing"
)

func TestLogger(t *testing.T) {
	cfg := &LoggerConfig{
		Service:   "myApp",
		Level:     DebugLevel,
		Dir:       "./log",
		Stdout:    true,
		ExtraKeys: []string{"key1", "key2"},
	}
	callerSkipOpt := WithZapOptions(zap.AddCallerSkip(3))
	// 手机号脱敏钩子函数
	var phoneDesensitizationHook = func(fields []zapcore.Field) {
		phoneRegex := regexp.MustCompile(`(\d{3})\d{4}(\d{4})`)
		for i := range fields {
			if fields[i].Type == zapcore.StringType {
				strValue := fields[i].String
				if phoneRegex.MatchString(strValue) {
					fields[i].String = phoneRegex.ReplaceAllString(strValue, `$1****$2`)
				}
			}
		}
	}
	phoneDesensitizationOpt := WithZapFieldHookFunc(phoneDesensitizationHook)
	if err := NewLogger(cfg, callerSkipOpt, phoneDesensitizationOpt); err != nil {
		assert.Nil(t, err)
	}
	defer Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key1", "value1")
	Info(ctx, "hello world")
	Infof(ctx, "hello %s", "world")
	Infow(ctx, "hello world", "key", "value")
	Error(ctx, "hello world")
	Errorf(ctx, "hello %s", "world")
	Errorw(ctx, "hello world", "key", "value")
	Infow(ctx, "phone info", "phone", "12312341234")

}

func TestExtraKeys(t *testing.T) {
	cfg := &LoggerConfig{
		Service:   "myApp",
		Level:     DebugLevel,
		Dir:       "./log",
		Stdout:    true,
		ExtraKeys: []string{"key1", "key2"},
	}
	opt := WithZapOptions(zap.AddCallerSkip(3))
	if err := NewLogger(cfg, opt); err != nil {
		assert.Nil(t, err)
	}
	defer Close()
	ctx := context.Background()
	ctx = context.WithValue(ctx, "key1", "value1")
	ctx = context.WithValue(ctx, "key2", "value2")
	Info(ctx, "hello world")
	Infof(ctx, "hello %s", "world")
	Infow(ctx, "hello world", "key", "value")
	Error(ctx, "hello world")
	Errorf(ctx, "hello %s", "world")
	ctx = context.WithValue(ctx, KeySkipLog, "")
	Errorw(ctx, "hello world", "key", "value11")
}
