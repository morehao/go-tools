package concq

import (
	"context"
	"fmt"
	"time"
)

func Example() {
	q := New(3, 10)

	for i := 0; i < 5; i++ {
		n := i
		err := q.Submit(func(ctx context.Context) error {
			fmt.Println("处理任务", n)
			time.Sleep(500 * time.Millisecond)
			return nil
		})
		if err != nil {
			fmt.Println("提交失败：", err)
		}
	}

	q.Shutdown()
	fmt.Println("所有任务完成")
}
