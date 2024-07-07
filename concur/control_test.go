package concur

import (
	"context"
	"fmt"
	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestControl_Run(t *testing.T) {
	defer glog.Close()
	mysqlClient, initDbErr := initMysqlClient()
	assert.Nil(t, initDbErr)
	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", "123456")
	// 实例化并发控制器，设置并发数为3
	cc := NewControl(3)
	var userIds []uint64
	for i := 0; i < 10000; i++ {
		userIds = append(userIds, uint64(i+1))
	}
	// 将ids分组，50一组
	var groupIds [][]uint64
	for i := 0; i < len(userIds); i += 50 {
		end := i + 50
		if end > len(userIds) {
			end = len(userIds)
		}
		groupIds = append(groupIds, userIds[i:end])
	}
	var result []interface{}
	for _, ids := range groupIds {
		tempIds := ids
		cc.Run(func() error {
			var userList []interface{}
			if err := mysqlClient.WithContext(ctx).Table("user").Where("id in ?", tempIds).Find(&userList).Error; err != nil {
				glog.Errorf(ctx, "query user err: %s, ids:%s", err, gutils.ToJson(tempIds))
				return err
			}
			result = append(result, userList...)
			return nil
		})
	}
	//
	// // 提交一些任务到并发控制器
	// for i := 0; i < 10; i++ {
	// 	count := i
	// 	cc.Run(func() error {
	// 		fmt.Printf("Processing task %d\n", count)
	// 		if count%2 == 0 { // 模拟一些任务失败
	// 			return fmt.Errorf("task %d failed", count)
	// 		}
	// 		return nil
	// 	})
	// }

	// 关闭并发控制器，等待所有任务完成
	cc.Close()

	// 获取并打印所有错误
	errors, failedCount := cc.Errors()
	fmt.Printf("Total failed tasks: %d\n", failedCount)
	for _, err := range errors {
		fmt.Println("Error:", err)
	}
}
