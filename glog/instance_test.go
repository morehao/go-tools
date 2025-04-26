package glog

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	ctx := context.Background()
	Debug(ctx, "message", "debug")
	Info(ctx, "message", "info")
	Warn(ctx, "message", "warn")
	Error(ctx, "message", "error")
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic after Panic")
		}
	}()
	Panic(ctx, "message", "fatal")
}

func TestLogLevels(t *testing.T) {
	ctx := context.Background()
	Debug(ctx, "message", "debug message")
	Info(ctx, "message", "info message")
	Warn(ctx, "message", "warn message")
	Error(ctx, "message", "error message")
}

func TestFileLogger(t *testing.T) {
	// 创建测试目录
	testDir := "./log"
	if err := os.MkdirAll(testDir, 0755); err != nil {
		t.Fatalf("failed to create test directory: %v", err)
	}
	defer os.RemoveAll(testDir)

	// 创建文件logger
	cfg := &LoggerConfig{
		module: "fs",
		Level:  InfoLevel,
		Type:   WriterFile,
		Dir:    testDir,
	}
	logger, err := newZapLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create file logger: %v", err)
	}

	// 设置新的默认logger
	oldLogger := defaultLogger
	SetDefaultLogger(logger)
	defer SetDefaultLogger(oldLogger)

	// 测试文件logger
	ctx := context.Background()
	Info(ctx, "message", "test file logger")

	// 检查日志文件是否创建
	files, err := filepath.Glob(filepath.Join(testDir, "*.log"))
	if err != nil {
		t.Fatalf("failed to glob log files: %v", err)
	}
	if len(files) == 0 {
		t.Error("no log file created")
	}
	t.Log(ToJsonString(files))
}

func TestClose(t *testing.T) {
	// 测试Close函数
	Close()

	// 测试Close后是否还能使用logger
	ctx := context.Background()
	defer func() {
		if r := recover(); r == nil {
			t.Error("expected panic after Close")
		}
	}()
	Info(ctx, "message", "this should panic")
}

// TestFieldHook 测试字段钩子函数
func TestFieldHook(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir, err := os.MkdirTemp("", "glog-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 设置测试配置
	config := &ServiceConfig{
		Service: "test",
		Modules: map[string]*LoggerConfig{
			"test": {
				module: "test",
				Level:  DebugLevel,
				Type:   WriterConsole,
				Dir:    tempDir,
			},
		},
	}

	// 记录钩子是否被调用
	var hookCalled bool
	var hookFields []Field

	var phoneDesensitizationHook = func(ctx context.Context, level Level, message string, fields ...Field) {
		hookCalled = true
		hookFields = fields
		phoneRegex := regexp.MustCompile(`(\d{3})\d{4}(\d{4})`)
		for i := range fields {
			if fields[i].Key == "phone" {
				strValue, ok := fields[i].Value.(string)
				if ok {
					if phoneRegex.MatchString(strValue) {
						fields[i].Value = phoneRegex.ReplaceAllString(strValue, `$1****$2`)
						t.Log("Phone number desensitized:", fields[i].Value)
					}
				}
			}
		}
	}

	var pwdDesensitizationHook = func(ctx context.Context, level Level, message string, fields ...Field) {
		// 处理消息中的密码
		if strings.Contains(message, "password") {
			re := regexp.MustCompile(`password=[^&\s]+`)
			message = re.ReplaceAllString(message, "password=***")
		}

		// 处理字段中的密码
		for i := range fields {
			if fields[i].Key == "password" {
				fields[i].Value = "***"
			}
		}
	}

	// 初始化日志器
	t.Log("Initializing logger with field hook")
	Init(config)
	AddHook(phoneDesensitizationHook)
	AddHook(pwdDesensitizationHook)

	// 测试电话号码脱敏
	ctx := context.Background()
	t.Log("Logging message with phone number")
	Infow(ctx, "test message", Field{Key: "phone", Value: "13812345678"})

	// 验证钩子是否被调用
	if !hookCalled {
		t.Error("Field hook was not called")
	}

	// 验证钩子接收到的字段
	if len(hookFields) == 0 {
		t.Error("No fields received by hook")
		return
	}

	t.Log("Hook fields:", hookFields)
	if hookFields[0].Key != "phone" || hookFields[0].Value != "138****5678" {
		t.Errorf("Unexpected field value: %v", hookFields[0])
	}

	// 测试密码脱敏
	t.Log("Logging message with password")
	Infow(ctx, "test message with password=123456", "password", "123456")
}

// TestContextLogger 测试上下文相关的logger操作，主要用于自定义日志组件等特殊场景
func TestContextLogger(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir, err := os.MkdirTemp("", "glog-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 设置测试配置
	config := &ServiceConfig{
		Service: "app",
		Modules: map[string]*LoggerConfig{
			"test": {
				module: "test",
				Level:  DebugLevel,
				Type:   WriterFile,
				Dir:    tempDir,
			},
		},
	}

	// 初始化日志器
	t.Log("Initializing logger")
	Init(config)

	// 创建测试上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "123456")
	ctx = context.WithValue(ctx, "user_id", "user123")

	// 获取上下文日志器
	logger := GetLogger(ctx)

	// 测试日志记录
	t.Log("Testing context logger")
	logger.Info(ctx, "message", "test context logger")

	// 测试带字段的日志
	logger.Infow(ctx, "test with fields", "key", "value")

	// 测试格式化日志
	logger.Infof(ctx, "test format: %s", "value", "test")
}

func TestSetDefaultLogger(t *testing.T) {
	// 创建一个新的logger
	cfg := &LoggerConfig{
		module: "test_default",
		Level:  DebugLevel,
		Type:   WriterConsole,
		Dir:    "./test_log",
	}
	logger, err := newZapLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// 保存旧的默认logger
	oldLogger := GetLogger(context.Background())

	// 设置新的默认logger
	SetDefaultLogger(logger)

	// 测试新的默认logger是否生效
	ctx := context.Background()
	actualLogger := GetLogger(ctx)
	if actualLogger != logger {
		t.Error("GetLogger should return the new default logger")
	}

	// 恢复旧的默认logger
	SetDefaultLogger(oldLogger)
}

// TestExtraKeys 测试从上下文中提取额外字段的功能
func TestExtraKeys(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir, err := os.MkdirTemp("", "glog-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 设置测试配置
	config := &ServiceConfig{
		Service: "test",
		Modules: map[string]*LoggerConfig{
			"test": {
				module:    "test",
				Level:     DebugLevel,
				Type:      WriterConsole,
				Dir:       tempDir,
				ExtraKeys: []string{"trace_id", "user_id", "request_id"},
			},
		},
	}

	// 初始化日志器
	t.Log("Initializing logger with extra keys")
	Init(config)

	// 获取模块级别的 logger
	logger := GetModuleLogger("test")

	// 创建带有额外字段的上下文
	ctx := context.Background()
	ctx = context.WithValue(ctx, "trace_id", "123456")
	ctx = context.WithValue(ctx, "user_id", "user123")
	ctx = context.WithValue(ctx, "request_id", "req789")
	// 添加一个不在 ExtraKeys 中的字段，用于测试过滤
	ctx = context.WithValue(ctx, "other_field", "should_not_appear")

	// 记录一条日志
	t.Log("Logging message with extra fields")
	logger.Infow(ctx, "test message with extra fields", "key", "value")

	// 同步日志
	Close()
}

func TestLogWithFields(t *testing.T) {
	ctx := context.Background()
	Infow(ctx, "info with fields", "key1", "value1", "key2", "value2")
	Errorw(ctx, "error with fields", "error", "something went wrong", "code", 500)
}

func TestLogFormat(t *testing.T) {
	ctx := context.Background()
	Debugf(ctx, "debug format: %s", "value", "test")
	Infof(ctx, "info format: %s", "value", "test")
	Warnf(ctx, "warn format: %s", "value", "test")
	Errorf(ctx, "error format: %s", "value", "test")
}

func TestHookFunctions(t *testing.T) {
	ctx := context.Background()

	// 添加手机号脱敏钩子
	AddHook(func(ctx context.Context, level Level, message string, fields ...Field) {
		for i, f := range fields {
			if f.Key == "phone" {
				phone := f.Value.(string)
				if len(phone) >= 11 {
					fields[i].Value = phone[:3] + "****" + phone[7:]
				}
			}
		}
	})

	// 添加密码脱敏钩子
	AddHook(func(ctx context.Context, level Level, message string, fields ...Field) {
		for i, f := range fields {
			if f.Key == "password" {
				fields[i].Value = "******"
			}
		}
	})

	Infow(ctx, "user login", "phone", "13812345678", "password", "123456")
}

func TestHookErrorHandling(t *testing.T) {
	ctx := context.Background()

	// 添加会panic的钩子
	AddHook(func(ctx context.Context, level Level, message string, fields ...Field) {
		panic("hook panic")
	})

	// 添加正常的钩子
	AddHook(func(ctx context.Context, level Level, message string, fields ...Field) {
		for i, f := range fields {
			if f.Key == "value" {
				fields[i].Value = "modified"
			}
		}
	})

	Info(ctx, "message", "test hook error handling", "value", "original")
}

func TestHookWithFields(t *testing.T) {
	ctx := context.Background()

	// 添加字段修改钩子
	AddHook(func(ctx context.Context, level Level, message string, fields ...Field) {
		for i, f := range fields {
			if f.Key == "value" {
				fields[i].Value = "modified by hook"
			}
		}
	})

	Infow(ctx, "test hook with fields", "value", "original value", "other", "unchanged")
}
