package conc

import (
	"context"
	"fmt"
	"testing"
	"time"
)

func TestControl_Run(t *testing.T) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", "123456")
	// 实例化并发控制器，设置并发数为3
	cc := NewControl(5)
	var userIds []uint64
	for i := 0; i < 100; i++ {
		userIds = append(userIds, uint64(i+1))
	}
	// 将ids分组，10个一组
	var groupIds [][]uint64
	for i := 0; i < len(userIds); i += 10 {
		end := i + 10
		if end > len(userIds) {
			end = len(userIds)
		}
		groupIds = append(groupIds, userIds[i:end])
	}
	for _, ids := range groupIds {
		tempIds := ids
		cc.Run(func(ctx context.Context) error {
			time.Sleep(time.Second * 3)
			fmt.Println("ids:", tempIds)
			return nil
		})
	}

	// 关闭并发控制器，等待所有任务完成
	cc.Close()

	// 获取并打印所有错误
	errors, failedCount := cc.Errors()
	fmt.Printf("Total failed tasks: %d\n", failedCount)
	for _, err := range errors {
		fmt.Println("Error:", err)
	}
}
