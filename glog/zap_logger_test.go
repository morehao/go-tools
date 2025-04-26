package glog

import (
	"context"
	"os"
	"testing"
)

// func TestLogger(t *testing.T) {
// 	cfg := &LoggerConfig{
// 		Service:   "myApp",
// 		Level:     DebugLevel,
// 		Dir:       "./log",
// 		Stdout:    true,
// 		ExtraKeys: []string{"key1", "key2"},
// 	}
// 	callerSkipOpt := WithZapOptions(zap.AddCallerSkip(3))
// 	// 手机号脱敏钩子函数
// 	var phoneDesensitizationHook = func(fields []zapcore.Field) {
// 		phoneRegex := regexp.MustCompile(`(\d{3})\d{4}(\d{4})`)
// 		for i := range fields {
// 			if fields[i].Type == zapcore.StringType {
// 				strValue := fields[i].String
// 				if phoneRegex.MatchString(strValue) {
// 					fields[i].String = phoneRegex.ReplaceAllString(strValue, `$1****$2`)
// 				}
// 			}
// 		}
// 	}
// 	phoneDesensitizationOpt := WithZapFieldHookFunc(phoneDesensitizationHook)
//
// 	var pwdDesensitizationHook = func(message string) string {
// 		// 只在消息中包含 "password" 关键字时进行脱敏处理
// 		if strings.Contains(message, "password") {
// 			// 匹配以 "password=" 开头的密码，并替换为脱敏的形式
// 			re := regexp.MustCompile(`password=[^&\s]+`)
// 			return re.ReplaceAllString(message, "password=***")
// 		}
// 		// 如果消息中不包含 "password" 关键字，则不进行处理
// 		return message
// 	}
// 	pwdDesensitizationOpt := WithMessageHookFunc(pwdDesensitizationHook)
// 	if err := NewLogger(cfg, callerSkipOpt, phoneDesensitizationOpt, pwdDesensitizationOpt); err != nil {
// 		assert.Nil(t, err)
// 	}
// 	defer Close()
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, "key1", "value1")
// 	Info(ctx, "hello world")
// 	Infof(ctx, "hello %s", "world")
// 	Infow(ctx, "hello world", "key", "value")
// 	Error(ctx, "hello world")
// 	Errorf(ctx, "hello %s", "world")
// 	Errorw(ctx, "hello world", "key", "value")
// 	Infow(ctx, "phone info", "phone", "12312341234")
// 	Info(ctx, "password=123456")
//
// }
//
// func TestExtraKeys(t *testing.T) {
// 	cfg := &LoggerConfig{
// 		Service:   "myApp",
// 		Level:     DebugLevel,
// 		Dir:       "./log",
// 		Stdout:    true,
// 		ExtraKeys: []string{"key1", "key2"},
// 	}
// 	opt := WithZapOptions(zap.AddCallerSkip(3))
// 	if err := NewLogger(cfg, opt); err != nil {
// 		assert.Nil(t, err)
// 	}
// 	defer Close()
// 	ctx := context.Background()
// 	ctx = context.WithValue(ctx, "key1", "value1")
// 	ctx = context.WithValue(ctx, "key2", "value2")
// 	Info(ctx, "hello world")
// 	Infof(ctx, "hello %s", "world")
// 	Infow(ctx, "hello world", "key", "value")
// 	Error(ctx, "hello world")
// 	Errorf(ctx, "hello %s", "world")
// 	ctx = context.WithValue(ctx, KeySkipLog, "")
// 	Errorw(ctx, "hello world", "key", "value11")
// }

func TestZapLoggerWithHooks(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir, err := os.MkdirTemp("", "glog-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 设置测试配置
	config := &LoggerConfig{
		service: "test",
		module:  "test",
		Level:   DebugLevel,
		Writer:  WriterFile,
		Dir:     tempDir,
	}

	// 记录钩子是否被调用
	var fieldHookCalled bool
	var messageHookCalled bool
	var hookFields []Field
	var hookMessage string

	// 创建日志器
	logger, err := newZapLogger(config,
		WithCallerSkip(1),
		WithFieldHookFunc(func(fields []Field) {
			fieldHookCalled = true
			hookFields = fields
		}),
		WithMessageHookFunc(func(msg string) string {
			messageHookCalled = true
			hookMessage = msg
			return "modified: " + msg
		}),
	)
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	// 记录一条日志
	ctx := context.Background()
	logger.Debugw(ctx, "test message", "key", "value")

	// 验证钩子是否被调用
	if !fieldHookCalled {
		t.Error("Field hook was not called")
	}
	if !messageHookCalled {
		t.Error("Message hook was not called")
	}

	// 验证钩子接收到的字段
	if len(hookFields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(hookFields))
	}
	if hookFields[0].Key != "key" || hookFields[0].Value != "value" {
		t.Errorf("Unexpected field value: %v", hookFields[0])
	}

	// 验证钩子接收到的消息
	if hookMessage != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", hookMessage)
	}
}
