package glog

import (
	"context"
	"os"
	"path/filepath"
	"regexp"
	"testing"
)

func TestDefaultLogger(t *testing.T) {
	// 测试默认logger是否正常初始化
	if defaultLogger == nil {
		t.Fatal("defaultLogger should not be nil")
	}

	// 测试默认logger的配置
	ctx := context.Background()
	Info(ctx, "test default logger")
}

func TestLogLevels(t *testing.T) {
	ctx := context.Background()

	// 测试各个日志级别
	Debug(ctx, "debug message")
	Info(ctx, "info message")
	Warn(ctx, "warn message")
	Error(ctx, "error message")

	// 测试格式化输出
	Debugf(ctx, "debug format: %s", "test")
	Infof(ctx, "info format: %s", "test")
	Warnf(ctx, "warn format: %s", "test")
	Errorf(ctx, "error format: %s", "test")

	// 测试带字段的日志
	Debugw(ctx, "debug with fields", "key", "value")
	Infow(ctx, "info with fields", "key", "value")
	Warnw(ctx, "warn with fields", "key", "value")
	Errorw(ctx, "error with fields", "key", "value")
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
		Service: "test",
		Module:  "fs",
		Level:   InfoLevel,
		Type:    WriterFile,
		Dir:     testDir,
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
	Info(ctx, "test file logger")

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
	Info(ctx, "this should panic")
}

// TestFieldHook 测试字段钩子函数
func TestFieldHook(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir, err := os.MkdirTemp("", "glog-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// defer os.RemoveAll(tempDir)

	// 设置测试配置
	config := &ServiceConfig{
		Service: "test",
		Modules: map[string]*LoggerConfig{
			"test": {
				Service: "test",
				Module:  "test",
				Level:   DebugLevel,
				Type:    WriterFile,
				Dir:     tempDir,
			},
		},
	}

	// 记录钩子是否被调用
	var hookCalled bool
	var hookFields []Field

	var phoneDesensitizationHook = func(fields []Field) {
		phoneRegex := regexp.MustCompile(`(\d{3})\d{4}(\d{4})`)
		for i := range fields {
			if fields[i].Key == "phone" {
				strValue, ok := fields[i].Value.(string)
				if ok {
					if phoneRegex.MatchString(strValue) {
						fields[i].Value = phoneRegex.ReplaceAllString(strValue, `$1****$2`)
					}
				}
			}
		}
	}

	// 初始化日志器
	Init(config, WithFieldHookFunc(phoneDesensitizationHook))

	// 记录一条日志
	ctx := context.Background()
	Debugw(ctx, "test message", "phone", "13812345678")

	// 验证钩子是否被调用
	if !hookCalled {
		t.Error("Field hook was not called")
	}

	// 验证钩子接收到的字段
	if len(hookFields) != 1 {
		t.Errorf("Expected 1 field, got %d", len(hookFields))
	}
	if hookFields[0].Key != "key" || hookFields[0].Value != "value" {
		t.Errorf("Unexpected field value: %v", hookFields[0])
	}
}

// TestMessageHook 测试消息钩子函数
func TestMessageHook(t *testing.T) {
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
				Service: "test",
				Module:  "test",
				Level:   DebugLevel,
				Type:    WriterFile,
				Dir:     tempDir,
			},
		},
	}

	// 记录钩子是否被调用
	var hookCalled bool
	var originalMessage string

	// 初始化日志器
	Init(config, WithMessageHookFunc(func(msg string) string {
		hookCalled = true
		originalMessage = msg
		return "modified: " + msg
	}))

	// 记录一条日志
	ctx := context.Background()
	Debug(ctx, "test message")

	// 验证钩子是否被调用
	if !hookCalled {
		t.Error("Message hook was not called")
	}

	// 验证钩子接收到的消息
	if originalMessage != "test message" {
		t.Errorf("Expected message 'test message', got '%s'", originalMessage)
	}
}

// TestContextLogger 测试上下文相关的logger操作，主要用于自定义日志组件等特殊场景
func TestContextLogger(t *testing.T) {
	// 创建一个新的logger
	cfg := &LoggerConfig{
		Module: "test",
		Level:  DebugLevel,
		Type:   WriterConsole,
		Dir:    "./test_log",
	}
	logger, err := newZapLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// 测试WithLogger
	ctx := context.Background()
	ctx = WithLogger(ctx, logger)

	// 测试GetLogger
	gotLogger := GetLogger(ctx)
	if gotLogger != logger {
		t.Error("GetLogger should return the logger from context")
	}

	// 测试日志输出
	gotLogger.Info(ctx, "test context logger")
}

func TestSetDefaultLogger(t *testing.T) {
	// 创建一个新的logger
	cfg := &LoggerConfig{
		Module: "test_default",
		Level:  DebugLevel,
		Type:   WriterConsole,
		Dir:    "./test_log",
	}
	logger, err := newZapLogger(cfg)
	if err != nil {
		t.Fatalf("failed to create logger: %v", err)
	}

	// 保存旧的默认logger
	oldLogger := defaultLogger

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

// TestModuleLogger 测试模块级别的logger
func TestModuleLogger(t *testing.T) {
	// 测试获取模块logger
	logger := GetModuleLogger("test_module")
	if logger == nil {
		t.Error("GetModuleLogger should not return nil")
	}

	// 测试日志输出
	ctx := context.Background()
	logger.Info(ctx, "test module logger")
}

// TestHookFunctions 测试钩子函数
func TestHookFunctions(t *testing.T) {
	ctx := context.Background()

	// 用于记录钩子函数的执行顺序
	var executionOrder []string

	// 添加多个钩子函数
	AddHook(func(ctx context.Context, level Level, msg string, fields ...Field) {
		executionOrder = append(executionOrder, "hook1")
	})

	AddHook(func(ctx context.Context, level Level, msg string, fields ...Field) {
		executionOrder = append(executionOrder, "hook2")
	})

	// 测试钩子函数的执行
	Info(ctx, "test hook functions")

	// 验证钩子函数的执行顺序
	if len(executionOrder) != 2 {
		t.Errorf("expected 2 hook executions, got %d", len(executionOrder))
	}
	if executionOrder[0] != "hook1" || executionOrder[1] != "hook2" {
		t.Error("hook functions executed in wrong order")
	}
}

// TestHookErrorHandling 测试钩子函数中的错误处理
func TestHookErrorHandling(t *testing.T) {
	ctx := context.Background()

	// 添加一个会panic的钩子函数
	AddHook(func(ctx context.Context, level Level, msg string, fields ...Field) {
		panic("hook panic")
	})

	// 添加一个正常的钩子函数
	var normalHookExecuted bool
	AddHook(func(ctx context.Context, level Level, msg string, fields ...Field) {
		normalHookExecuted = true
	})

	// 测试钩子函数的错误处理
	Info(ctx, "test hook error handling")

	// 验证正常的钩子函数仍然执行
	if !normalHookExecuted {
		t.Error("normal hook should still execute after panic")
	}
}

// TestHookWithFields 测试带字段的钩子函数
func TestHookWithFields(t *testing.T) {
	ctx := context.Background()

	// 用于记录接收到的字段
	var receivedFields []Field

	// 添加一个处理字段的钩子函数
	AddHook(func(ctx context.Context, level Level, msg string, fields ...Field) {
		receivedFields = fields
	})

	// 测试带字段的日志
	Infow(ctx, "test hook with fields", "key1", "value1", "key2", "value2")

	// 验证字段是否正确传递
	if len(receivedFields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(receivedFields))
	}
}

// TestHookLevelFilter 测试不同日志级别的钩子函数
func TestHookLevelFilter(t *testing.T) {
	ctx := context.Background()

	// 用于记录钩子函数的执行次数
	var hookExecutions int

	// 添加一个钩子函数
	AddHook(func(ctx context.Context, level Level, msg string, fields ...Field) {
		hookExecutions++
	})

	// 测试不同日志级别
	Debug(ctx, "debug message")
	Info(ctx, "info message")
	Warn(ctx, "warn message")
	Error(ctx, "error message")

	// 验证钩子函数对每个日志级别都执行
	if hookExecutions != 4 {
		t.Errorf("expected 4 hook executions, got %d", hookExecutions)
	}
}

// TestHookContext 测试钩子函数中的上下文传递
func TestHookContext(t *testing.T) {
	// 创建一个带有特定值的上下文
	ctx := context.WithValue(context.Background(), "test_key", "test_value")

	// 用于验证上下文传递
	var contextValue string

	// 添加一个检查上下文的钩子函数
	AddHook(func(ctx context.Context, level Level, msg string, fields ...Field) {
		if val, ok := ctx.Value("test_key").(string); ok {
			contextValue = val
		}
	})

	// 测试日志输出
	Info(ctx, "test hook context")

	// 验证上下文值是否正确传递
	if contextValue != "test_value" {
		t.Error("context value not correctly passed to hook")
	}
}
