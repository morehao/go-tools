package glog

import (
	"context"
	"os"
	"path/filepath"
	"testing"
	"time"
)

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
		Init(config)

		// 验证默认 logger 是否创建成功
		logger := GetLogger(context.Background())
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
		Init(config)

		// 验证各个模块的 logger
		module1Logger := GetModuleLogger("module1")
		if module1Logger == nil {
			t.Error("Module1 logger not initialized")
		}

		module2Logger := GetModuleLogger("module2")
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
		Init(config)

		// 验证 console logger
		logger := GetModuleLogger("console")
		if logger == nil {
			t.Error("Console logger not initialized")
		}

		// 写入日志（这里主要测试不会panic）
		ctx := context.Background()
		logger.Debug(ctx, "debug to console")
		logger.Info(ctx, "info to console")
	})

	t.Run("TestLoggerWithHooks", func(t *testing.T) {
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

		// 添加钩子函数
		var hookCalled bool

		// 初始化日志系统
		Init(config)

		// 写入日志
		ctx := context.Background()
		Info(ctx, "test hook")

		// 验证钩子是否被调用
		if !hookCalled {
			t.Error("Hook function not called")
		}
	})

	t.Run("TestLoggerWithExtraKeys", func(t *testing.T) {
		config := &LogConfig{
			Service: "test-service",
			Modules: map[string]*ModuleLoggerConfig{
				"default": {
					Level:     InfoLevel,
					Writer:    WriterFile,
					Dir:       tempDir,
					ExtraKeys: []string{"trace_id", "user_id"},
				},
			},
		}

		// 初始化日志系统
		Init(config)

		// 创建带有额外字段的上下文
		ctx := context.Background()
		ctx = context.WithValue(ctx, "trace_id", "123456")
		ctx = context.WithValue(ctx, "user_id", "user123")

		// 写入日志
		Info(ctx, "test extra keys")

		// 验证日志文件
		expectedDir := filepath.Join(tempDir, time.Now().Format("20060102"))
		expectedFile := filepath.Join(expectedDir, "test-service_full.log")
		if !fileExists(expectedFile) {
			t.Errorf("Log file not created: %s", expectedFile)
		}
	})
}
