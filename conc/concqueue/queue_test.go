package concqueue

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync/atomic"
	"testing"
	"time"
)

// 自定义 Logger
type testLogger struct {
	logger *log.Logger
	errors []string // 用于记录错误信息
}

func (l *testLogger) Errorf(ctx context.Context, format string, args ...any) {
	// 在这里可以自定义日志的格式和输出方式
	msg := fmt.Sprintf(format, args...)
	l.logger.Printf("[ERROR] " + msg)
	l.errors = append(l.errors, msg) // 记录错误信息

	// 打印 Context 中的信息
	for _, key := range []interface{}{"request_id", "user_id"} {
		if value := ctx.Value(key); value != nil {
			l.logger.Printf("  %v: %v", key, value)
		}
	}
}

func Test_Example(t *testing.T) {
	Handler := func(ctx context.Context, input int) (string, error) {
		fmt.Println(input, "处理中...")
		time.Sleep(time.Millisecond * 100) // 缩短睡眠时间，加快测试
		if input%2 == 0 {
			return "", fmt.Errorf("输入参数为偶数")
		}
		res := fmt.Sprintf("处理结果：%d", input)
		return res, nil
	}

	start := time.Now()

	// 创建自定义 Logger
	customLogger := &testLogger{
		logger: log.New(os.Stdout, "test: ", log.LstdFlags),
		errors: []string{},
	}

	// 创建队列，使用3个worker，队列大小为10
	q := New(5, 10, WithLogger(customLogger))
	taskCount := 5
	res := make([]string, taskCount)
	errs := make([]error, taskCount)

	// 提交任务
	for i := 0; i < taskCount; i++ {
		n := i // 传递索引值给任务，避免并发时数据竞争
		q.Submit(func(ctx context.Context) error {
			currRes, err := Handler(ctx, n)
			if err != nil {
				errs[n] = err
			} else {
				res[n] = currRes
			}
			return err
		})
	}

	// 等待任务完成并关闭队列
	errCnt := q.Shutdown()

	// 打印并检查错误数量
	t.Logf("错误数量：%d", errCnt)

	// 验证结果
	for i := 0; i < taskCount; i++ {
		if i%2 == 0 {
			if errs[i] == nil {
				t.Errorf("任务 %d 应该返回错误，但没有返回错误", i)
			}
		} else {
			if errs[i] != nil {
				t.Errorf("任务 %d 不应返回错误，但返回了错误: %v", i, errs[i])
			} else if res[i] != fmt.Sprintf("处理结果：%d", i) {
				t.Errorf("任务 %d 结果不正确，期望 %s，但实际是 %s", i, fmt.Sprintf("处理结果：%d", i), res[i])
			}
		}
	}
	// 断言错误数量正确
	var expectedErrCount int32 = 3 // 偶数任务会返回错误
	if int(errCnt) != int(expectedErrCount) {
		t.Errorf("错误数量不符，期望 %d，但实际是 %d", expectedErrCount, errCnt)
	}

	// 验证自定义 Logger 是否记录了错误
	if len(customLogger.errors) != int(expectedErrCount) {
		t.Errorf("自定义 Logger 记录的错误数量不符，期望 %d，但实际是 %d", expectedErrCount, len(customLogger.errors))
	}

	t.Log("所有任务完成，耗时：", time.Since(start))
	t.Log("所有任务完成")
}

func TestContextKeys(t *testing.T) {
	// 创建自定义 Logger
	customLogger := &testLogger{
		logger: log.New(os.Stdout, "test: ", log.LstdFlags),
		errors: []string{},
	}

	// 创建队列，并设置 Context Keys
	q := New(
		1, 1,
		WithLogger(customLogger),
		WithContextKeys("request_id", "user_id"),
	)

	// 创建带有 Context Key 的 Context
	ctx := context.WithValue(context.Background(), "request_id", "req-123")
	ctx = context.WithValue(ctx, "user_id", "user-456")

	// 提交一个会产生错误的 Task
	q.Submit(func(ctx context.Context) error {
		customLogger.Errorf(ctx, "This is an error message with context.")
		return fmt.Errorf("test error")
	})

	// 关闭队列
	q.Shutdown()

	// 验证自定义 Logger 是否记录了错误
	if len(customLogger.errors) != 2 { // 修改为期望 2 条错误
		t.Fatalf("Expected 2 error messages, got %d", len(customLogger.errors))
	}

	// 验证第一条消息
	expectedLog1 := "This is an error message with context."
	if customLogger.errors[0] != expectedLog1 {
		t.Errorf("Expected first log message '%s', got '%s'", expectedLog1, customLogger.errors[0])
	}

	// 验证第二条消息
	expectedLog2 := "task done, err: test error"
	if customLogger.errors[1] != expectedLog2 {
		t.Errorf("Expected second log message '%s', got '%s'", expectedLog2, customLogger.errors[1])
	}
}

func TestNoErrors(t *testing.T) {
	start := time.Now()
	// 创建自定义 Logger
	customLogger := &testLogger{
		logger: log.New(os.Stdout, "test: ", log.LstdFlags),
		errors: []string{},
	}

	// 创建队列，使用3个worker，队列大小为10
	q := New(5, 10, WithLogger(customLogger))
	taskCount := 5
	res := make([]string, taskCount)
	errs := make([]error, taskCount)

	// 提交任务
	for i := 0; i < taskCount; i++ {
		n := i // 传递索引值给任务，避免并发时数据竞争
		q.Submit(func(ctx context.Context) error {
			currRes, err := func(ctx context.Context, input int) (string, error) {
				fmt.Println(input, "处理中...")
				time.Sleep(time.Millisecond * 100) // 缩短睡眠时间，加快测试
				res := fmt.Sprintf("处理结果：%d", input)
				return res, nil
			}(ctx, n)
			if err != nil {
				errs[n] = err
			} else {
				res[n] = currRes
			}
			return nil
		})
	}

	// 等待任务完成并关闭队列
	errCnt := q.Shutdown()

	// 打印并检查错误数量
	t.Logf("错误数量：%d", errCnt)

	// 验证结果
	for i := 0; i < taskCount; i++ {
		if errs[i] != nil {
			t.Errorf("任务 %d 不应返回错误，但返回了错误: %v", i, errs[i])
		} else if res[i] != fmt.Sprintf("处理结果：%d", i) {
			t.Errorf("任务 %d 结果不正确，期望 %s，但实际是 %s", i, fmt.Sprintf("处理结果：%d", i), res[i])
		}
	}
	// 断言错误数量正确
	var expectedErrCount int32 = 0
	if int(errCnt) != int(expectedErrCount) {
		t.Errorf("错误数量不符，期望 %d，但实际是 %d", expectedErrCount, errCnt)
	}

	// 验证自定义 Logger 是否记录了错误
	if len(customLogger.errors) != int(expectedErrCount) {
		t.Errorf("自定义 Logger 记录的错误数量不符，期望 %d，但实际是 %d", expectedErrCount, len(customLogger.errors))
	}

	t.Log("所有任务完成，耗时：", time.Since(start))
	t.Log("所有任务完成")
}

func TestPanic(t *testing.T) {
	// 创建自定义 Logger
	customLogger := &testLogger{
		logger: log.New(os.Stdout, "test: ", log.LstdFlags),
		errors: []string{},
	}

	// 创建队列，使用3个worker，队列大小为10
	var panicCount int32
	q := New(5, 10, WithLogger(customLogger), WithPanicHandler(func(r interface{}) {
		atomic.AddInt32(&panicCount, 1)
		fmt.Printf("捕获到panic: %v\n", r)
	}))
	taskCount := 5
	res := make([]string, taskCount)
	errs := make([]error, taskCount)

	// 提交任务
	for i := 0; i < taskCount; i++ {
		n := i // 传递索引值给任务，避免并发时数据竞争
		q.Submit(func(ctx context.Context) error {
			currRes, err := func(ctx context.Context, input int) (string, error) {
				fmt.Println(input, "处理中...")
				time.Sleep(time.Millisecond * 100) // 缩短睡眠时间，加快测试
				if input == 2 {
					panic("模拟panic")
				}
				res := fmt.Sprintf("处理结果：%d", input)
				return res, nil
			}(ctx, n)
			if err != nil {
				errs[n] = err
			} else {
				res[n] = currRes
			}
			return nil
		})
	}

	// 等待任务完成并关闭队列
	errCnt := q.Shutdown()

	// 打印并检查错误数量
	t.Logf("错误数量：%d", errCnt)

	// 验证结果
	for i := 0; i < taskCount; i++ {
		if i == 2 {
			continue
		}
		if errs[i] != nil {
			t.Errorf("任务 %d 不应返回错误，但返回了错误: %v", i, errs[i])
		} else if res[i] != fmt.Sprintf("处理结果：%d", i) {
			t.Errorf("任务 %d 结果不正确，期望 %s，但实际是 %s", i, fmt.Sprintf("处理结果：%d", i), res[i])
		}
	}
	// 断言错误数量正确
	var expectedErrCount int32 = 0
	if int(errCnt) != int(expectedErrCount) {
		t.Errorf("错误数量不符，期望 %d，但实际是 %d", expectedErrCount, errCnt)
	}

	// 验证自定义 Logger 是否记录了错误
	if len(customLogger.errors) != int(expectedErrCount) {
		t.Errorf("自定义 Logger 记录的错误数量不符，期望 %d，但实际是 %d", expectedErrCount, len(customLogger.errors))
	}

	if panicCount != 1 {
		t.Errorf("Panic 数量不符，期望 %d，但实际是 %d", 1, panicCount)
	}

	t.Log("所有任务完成")
}

func TestTimeout(t *testing.T) {
	// 创建自定义 Logger
	customLogger := &testLogger{
		logger: log.New(os.Stdout, "test: ", log.LstdFlags),
		errors: []string{},
	}

	// 创建队列，使用1个worker，队列大小为1
	q := New(1, 1, WithSubmitTimeout(time.Millisecond*200), WithLogger(customLogger))

	// 提交任务
	q.Submit(func(ctx context.Context) error {
		fmt.Println("任务 0 开始执行...")
		time.Sleep(time.Millisecond * 500) // 模拟任务执行时间超过超时
		fmt.Println("任务 0 执行完毕")
		return nil
	})

	// 提交另一个任务
	q.Submit(func(ctx context.Context) error {
		fmt.Println("任务 1 开始执行...")
		time.Sleep(time.Millisecond * 500) // 模拟任务执行时间超过超时
		fmt.Println("任务 1 执行完毕")
		return nil
	})

	// 等待任务完成并关闭队列
	errCnt := q.Shutdown()

	// 打印并检查错误数量
	t.Logf("错误数量：%d", errCnt)

	// 检查是否触发了超时
	if len(customLogger.errors) == 0 {
		t.Error("期望日志中记录 submit timeout，但实际没有")
	}

	t.Log("测试完成")
}
