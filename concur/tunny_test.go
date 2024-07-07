package concur

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/morehao/go-tools/dbClient"
	"github.com/morehao/go-tools/glog"
	"github.com/morehao/go-tools/gutils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"sync"
)

func initMysqlClient() (*gorm.DB, error) {
	if err := glog.InitZapLogger(&glog.LoggerConfig{
		Service:   "test",
		Level:     glog.DebugLevel,
		Dir:       "./log",
		Stdout:    true,
		ExtraKeys: []string{"requestId"},
	}); err != nil {
		return nil, err
	}

	cfg := dbClient.MysqlConfig{
		Service:  "test",
		Addr:     "127.0.0.1:3306",
		Database: "demo",
		User:     "root",
		Password: "123456",
	}
	return dbClient.InitMysql(cfg)
}

func TestNewControl(t *testing.T) {
	defer glog.Close()
	mysqlClient, initDbErr := initMysqlClient()
	assert.Nil(t, initDbErr)

	fn := func(ctx context.Context, payload interface{}) interface{} {
		ids, ok := payload.([]uint64)
		if !ok {
			return errors.New("invalid payload type")
		}

		var userList []interface{}
		if err := mysqlClient.WithContext(ctx).Table("user").Where("id in ?", ids).Find(&userList).Error; err != nil {
			glog.Errorf(ctx, "query user err: %s, ids:%s", err, gutils.ToJson(ids))
			return nil
		}
		time.Sleep(100 * time.Millisecond) // 模拟长时间任务
		return userList
	}

	ctx := context.Background()
	ctx = context.WithValue(ctx, "requestId", "123456")
	ctrl := NewTunnyCtl(10, fn)

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
	var wg sync.WaitGroup
	for _, ids := range groupIds {
		wg.Add(1)
		go func(tempIds []uint64) {
			defer wg.Done()
			result = append(result, ctrl.Run(ctx, tempIds))
		}(ids)
	}
	wg.Wait()

	t.Log(gutils.ToJson(result))

	assert.Equal(t, int64(0), ctrl.GetErrCnt(), "Error count should be 0")
	ctrl.Close()
}
