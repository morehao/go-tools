package glog

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestInit(t *testing.T) {
	// 创建测试目录
	tempDir := "log/glog-test"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// defer os.RemoveAll(tempDir)

	t.Run("TestBasicInit", func(t *testing.T) {
		config := &LogConfig{
			Service: "test-service",
			Modules: map[string]*ModuleLoggerConfig{
				"default": {
					Level:  InfoLevel,
					Writer: WriterFile,
					Dir:    tempDir,
				},
			},
		}

		// 初始化日志系统
		InitLogger(config)

		// 验证默认 logger 是否创建成功
		logger := getLoggerFromCtx(context.Background())
		if logger == nil {
			t.Error("Default logger not initialized")
		}

		// 写入一条日志
		logger.Info(context.Background(), "test message")

		// 验证日志文件是否创建
		expectedDir := filepath.Join(tempDir, time.Now().Format("20060102"))
		expectedFile := filepath.Join(expectedDir, "test-service_full.log")
		if !fileExists(expectedFile) {
			t.Errorf("Log file not created: %s", expectedFile)
		}
	})

	t.Run("TestMultiModuleInit", func(t *testing.T) {
		config := &LogConfig{
			Service: "test-service",
			Modules: map[string]*ModuleLoggerConfig{
				"module1": {
					Level:  DebugLevel,
					Writer: WriterFile,
					Dir:    tempDir,
				},
				"module2": {
					Level:  InfoLevel,
					Writer: WriterFile,
					Dir:    tempDir,
				},
			},
		}

		// 初始化日志系统
		InitLogger(config)

		// 验证各个模块的 logger
		module1Logger, getModule1LoggerErr := GetModuleLogger("module1")
		assert.Nil(t, getModule1LoggerErr)
		if module1Logger == nil {
			t.Error("Module1 logger not initialized")
		}

		module2Logger, getModule2LoggerErr := GetModuleLogger("module2")
		assert.Nil(t, getModule2LoggerErr)
		if module2Logger == nil {
			t.Error("Module2 logger not initialized")
		}

		// 写入日志
		ctx := context.Background()
		module1Logger.Debug(ctx, "debug message")
		module2Logger.Info(ctx, "info message")

		// 验证日志文件
		expectedDir := filepath.Join(tempDir, time.Now().Format("20060102"))
		module1File := filepath.Join(expectedDir, "test-service_full.log")
		module2File := filepath.Join(expectedDir, "test-service_full.log")

		if !fileExists(module1File) {
			t.Errorf("Module1 log file not created: %s", module1File)
		}
		if !fileExists(module2File) {
			t.Errorf("Module2 log file not created: %s", module2File)
		}
	})

	t.Run("TestConsoleLogger", func(t *testing.T) {
		config := &LogConfig{
			Service: "test-service",
			Modules: map[string]*ModuleLoggerConfig{
				"console": {
					Level:  DebugLevel,
					Writer: WriterConsole,
					Dir:    tempDir,
				},
			},
		}

		// 初始化日志系统
		InitLogger(config)

		// 验证 console logger
		logger, getLoggerErr := GetModuleLogger("console")
		assert.Nil(t, getLoggerErr)
		if logger == nil {
			t.Error("Console logger not initialized")
		}

		// 写入日志（这里主要测试不会panic）
		ctx := context.Background()
		logger.Debug(ctx, "debug to console")
		logger.Info(ctx, "info to console")
	})
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
func TestHook(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir := "log/glog-test"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 设置测试配置
	config := &LogConfig{
		Service: "test",
		Modules: map[string]*ModuleLoggerConfig{
			"default": {
				Level:  DebugLevel,
				Writer: WriterConsole,
				Dir:    tempDir,
			},
		},
	}

	var phoneDesensitizationHook = func(fields []Field) {
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

	var pwdDesensitizationHook = func(message string) string {
		// 处理消息中的密码
		if strings.Contains(message, "password") {
			re := regexp.MustCompile(`password=[^&\s]+`)
			return re.ReplaceAllString(message, "password=***")
		}
		return message
	}

	// 初始化日志器
	t.Log("Initializing logger with field hook")
	InitLogger(config, WithFieldHookFunc(phoneDesensitizationHook), WithMessageHookFunc(pwdDesensitizationHook))

	// 测试电话号码脱敏
	ctx := context.Background()
	t.Log("Logging message with phone number")
	Infow(ctx, "test message", "phone", "13812345678")

	// 测试密码脱敏
	t.Log("Logging message with password")
	Info(ctx, "test message with password=123456")
}

// TestExtraKeys 测试从上下文中提取额外字段的功能
func TestExtraKeys(t *testing.T) {
	// 创建一个临时目录用于测试
	tempDir := "log/glog-test"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// 设置测试配置
	config := &LogConfig{
		Service: "test",
		Modules: map[string]*ModuleLoggerConfig{
			"test": {
				module:    "test",
				Level:     DebugLevel,
				Writer:    WriterConsole,
				Dir:       tempDir,
				ExtraKeys: []string{"trace_id", "user_id", "request_id"},
			},
		},
	}

	// 初始化日志器
	t.Log("Initializing logger with extra keys")
	InitLogger(config)

	// 获取模块级别的 logger
	logger, getLoggerErr := GetModuleLogger("test")
	if getLoggerErr != nil {
		t.Fatalf("failed to get logger: %v", getLoggerErr)
	}

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
	Debugf(ctx, "debug format: %s", "value")
	Infof(ctx, "info format: %s", "value")
	Warnf(ctx, "warn format: %s", "value")
	Errorf(ctx, "error format: %s", "value")
}

func TestRotateUnit(t *testing.T) {
	// 创建临时目录
	tempDir := "log/glog-test"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	// defer os.RemoveAll(tempDir)

	// 使用固定的时间戳
	now := time.Now()

	// 测试按天切割
	t.Run("TestRotateUnitDay", func(t *testing.T) {
		config := &LogConfig{
			Service: "test",
			Modules: map[string]*ModuleLoggerConfig{
				"test": {
					Level:      InfoLevel,
					Writer:     WriterFile,
					Dir:        tempDir,
					RotateUnit: RotateUnitDay,
				},
			},
		}

		// 初始化日志器
		InitLogger(config)

		// 记录日志
		ctx := context.Background()
		Info(ctx, "test message")

		// 验证日志文件是否存在
		expectedDir := filepath.Join(tempDir, now.Format("20060102"))
		expectedFile := filepath.Join(expectedDir, "test_full.log")
		if !fileExists(expectedFile) {
			t.Errorf("Expected log file %s does not exist", expectedFile)
		}
	})

	// 测试按小时切割
	t.Run("TestRotateUnitHour", func(t *testing.T) {
		config := &LogConfig{
			Service: "test",
			Modules: map[string]*ModuleLoggerConfig{
				"test": {
					Level:      InfoLevel,
					Writer:     WriterFile,
					Dir:        tempDir,
					RotateUnit: RotateUnitHour,
				},
			},
		}

		// 初始化日志器
		InitLogger(config)

		// 记录日志
		ctx := context.Background()
		Info(ctx, "test message")

		// 验证日志文件是否存在
		expectedDir := filepath.Join(tempDir, now.Format("20060102"))
		expectedFile := filepath.Join(expectedDir, fmt.Sprintf("test_full_%s.log", now.Format("15")))
		if !fileExists(expectedFile) {
			t.Errorf("Expected log file %s does not exist", expectedFile)
		}
	})

	// 测试默认值
	t.Run("TestDefaultRotateUnit", func(t *testing.T) {
		config := &LogConfig{
			Service: "app",
			Modules: map[string]*ModuleLoggerConfig{
				"default": {
					Level:  InfoLevel,
					Writer: WriterFile,
					Dir:    tempDir,
				},
			},
		}

		// 初始化日志器
		InitLogger(config)

		// 记录日志
		ctx := context.Background()
		Info(ctx, "test message")

		// 验证日志文件是否存在
		rootDir, _ := os.Getwd()
		expectedDir := filepath.Join(rootDir, tempDir, now.Format("20060102"))
		expectedFile := filepath.Join(expectedDir, "app_full.log")
		if !fileExists(expectedFile) {
			t.Errorf("Expected log file %s does not exist", expectedFile)
		}
	})
}
