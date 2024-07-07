package concur

import (
	"sync"
)

// Control 控制并发执行任务的控制器
type Control struct {
	wg        sync.WaitGroup
	workerNum int           // 并发数
	mu        sync.Mutex    // 用于保护下面的 errors 切片
	errors    []error       // 存储任务执行过程中的错误
	sem       chan struct{} // 用于控制并发数的信号量
}

// NewControl 创建一个新的并发控制器，workerNum 指定并发执行的工作数
func NewControl(workerNum int) *Control {
	return &Control{
		workerNum: workerNum,
		sem:       make(chan struct{}, workerNum), // 初始化信号量
	}
}

// Run 执行传入的方法，控制并发执行的数量
func (cc *Control) Run(task func() error) {
	cc.sem <- struct{}{} // 获取信号量
	cc.wg.Add(1)         // 增加等待计数
	go func() {
		defer cc.wg.Done()          // 完成后减少等待计数
		defer func() { <-cc.sem }() // 释放信号量
		if err := task(); err != nil {
			cc.mu.Lock()
			cc.errors = append(cc.errors, err)
			cc.mu.Unlock()
		}
	}()
}

// Close 关闭控制器，等待所有任务完成
func (cc *Control) Close() {
	cc.wg.Wait() // 等待所有任务完成
}

// Errors 返回执行过程中发生的所有错误及其数量
func (cc *Control) Errors() ([]error, int) {
	cc.mu.Lock()
	defer cc.mu.Unlock()
	return cc.errors, len(cc.errors)
}
