package concq

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func Test_Example(t *testing.T) {
	Handler := func(ctx context.Context, input int) (string, error) {
		fmt.Println(input, "处理中...")
		time.Sleep(time.Second * 2)
		if input%2 == 0 {
			return "", fmt.Errorf("输入参数为偶数")
		}
		res := fmt.Sprintf("处理结果：%d", input)
		return res, nil
	}

	start := time.Now()

	// 创建队列，使用3个worker，队列大小为10
	q := New(5, 10)
	taskCount := 5
	res := make([]string, taskCount)
	errs := make([]error, taskCount)

	// 提交任务
	for i := 0; i < taskCount; i++ {
		n := i // 传递索引值给任务，避免并发时数据竞争
		err := q.Submit(func(ctx context.Context) error {
			currRes, err := Handler(ctx, n)
			if err != nil {
				errs[n] = err
			} else {
				res[n] = currRes
			}
			return err
		})
		if err != nil {
			t.Errorf("提交任务失败：%v", err) // 使用 t.Errorf 来报告测试失败
		}
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
	expectedErrCount := 3 // 偶数任务会返回错误
	if errCnt != expectedErrCount {
		t.Errorf("错误数量不符，期望 %d，但实际是 %d", expectedErrCount, errCnt)
	}
	t.Log("所有任务完成，耗时：", time.Since(start))
	t.Log("所有任务完成")
}
