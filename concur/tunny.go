package concur

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"

	"github.com/Jeffail/tunny"
)

// TunnyCtrl 结构体定义
type TunnyCtrl struct {
	concurNum int
	pool      *tunny.Pool
	errCnt    int64
	errList   []error
	lock      sync.Mutex
}

// TunnyWorkerFn 函数类型定义
type TunnyWorkerFn func(context.Context, interface{}) interface{}

// NewTunnyCtl 创建一个新的 TunnyCtrl 实例
func NewTunnyCtl(concurNum int, fn TunnyWorkerFn) *TunnyCtrl {
	ctrl := &TunnyCtrl{
		concurNum: concurNum,
		pool: tunny.NewFunc(concurNum, func(payload interface{}) interface{} {
			p := payload.(struct {
				Ctx     context.Context
				Payload interface{}
			})
			return fn(p.Ctx, p.Payload)
		}),
	}
	return ctrl
}

// Run 提交一个任务到池中执行
func (ctl *TunnyCtrl) Run(ctx context.Context, payload interface{}) interface{} {
	result := ctl.pool.Process(struct {
		Ctx     context.Context
		Payload interface{}
	}{Ctx: ctx, Payload: payload})
	if result == nil {
		atomic.AddInt64(&ctl.errCnt, 1)
		ctl.lock.Lock()
		ctl.errList = append(ctl.errList, errors.New("process error"))
		ctl.lock.Unlock()
	}
	return result
}

// GetErrCnt 获取错误计数
func (ctl *TunnyCtrl) GetErrCnt() int64 {
	return atomic.LoadInt64(&ctl.errCnt)
}

// Close 关闭池
func (ctl *TunnyCtrl) Close() {
	ctl.pool.Close()
}
