package conc

import (
	"context"
	"sync"
)

// Control 控制并发执行任务的控制器
type Control struct {
	wg        sync.WaitGroup
	workerNum int           // 并发数
	mu        sync.Mutex    // 用于保护下面的 errors 切片
	errors    []error       // 存储任务执行过程中的错误
	sem       chan struct{} // 用于控制并发数的信号量
	ctx       context.Context
	cancel    context.CancelFunc
}

// NewControl 创建一个新的并发控制器，workerNum 指定并发执行的工作数
func NewControl(workerNum int) *Control {
	ctx, cancel := context.WithCancel(context.Background())
	return &Control{
		workerNum: workerNum,
		sem:       make(chan struct{}, workerNum), // 初始化信号量
		ctx:       ctx,
		cancel:    cancel,
	}
}

// Run 执行传入的方法，控制并发执行的数量
func (ctrl *Control) Run(task func(context.Context) error) {
	select {
	case ctrl.sem <- struct{}{}: // 获取信号量
	case <-ctrl.ctx.Done(): // 检查是否已经调用了cancel
		return
	}

	ctrl.wg.Add(1) // 增加等待计数
	go func() {
		defer ctrl.wg.Done()          // 完成后减少等待计数
		defer func() { <-ctrl.sem }() // 释放信号量

		// 执行任务并处理错误
		if err := task(ctrl.ctx); err != nil {
			ctrl.mu.Lock()
			ctrl.errors = append(ctrl.errors, err)
			ctrl.mu.Unlock()
		}
	}()
}

// Close 关闭控制器，等待所有任务完成
func (ctrl *Control) Close() {
	ctrl.cancel()  // 取消所有剩余的任务
	ctrl.wg.Wait() // 等待所有任务完成
}

// Errors 返回执行过程中发生的所有错误及其数量
func (ctrl *Control) Errors() ([]error, int) {
	ctrl.mu.Lock()
	defer ctrl.mu.Unlock()
	// 返回错误切片的副本
	errorsCopy := make([]error, len(ctrl.errors))
	copy(errorsCopy, ctrl.errors)
	return errorsCopy, len(errorsCopy)
}
