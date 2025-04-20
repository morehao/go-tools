package concpool

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
)

// Task 表示要执行的任务类型，加入了 context.Context
type Task func(ctx context.Context) error

// Pool 定义了对外公开的接口
type Pool interface {
	Submit(task Task) // 提交任务
	Shutdown() int32  // 关闭池，返回失败的任务数
}

// pool 是 Pool 接口的实现
// 现在 pool 是包内私有的结构体
type pool struct {
	taskQueue   chan Task      // 任务队列
	workerCount int            // worker 数量
	wg          sync.WaitGroup // 等待任务完成
	ctx         context.Context
	cancel      context.CancelFunc
	errCount    int32           // 使用原子操作记录失败的任务数
	onErr       func(err error) // 处理任务失败时的回调函数
	closed      int32           // 使用原子操作处理池是否已关闭的状态
}

// New 创建一个新的 pool，并自动启动
func New(workerCount int, queueSize int, options ...Option) Pool {
	ctx, cancel := context.WithCancel(context.Background())

	p := &pool{
		taskQueue:   make(chan Task, queueSize),
		workerCount: workerCount,
		ctx:         ctx,
		cancel:      cancel,
		onErr: func(err error) {
			fmt.Println(err)
		},
	}

	// 应用 options 配置
	for _, option := range options {
		option(p)
	}

	// 自动启动 pool
	p.start()

	return p
}

// start 启动 pool，改为私有方法
func (p *pool) start() {
	// 使用原子操作检查是否已经关闭
	if atomic.LoadInt32(&p.closed) == 1 {
		return
	}

	for i := 0; i < p.workerCount; i++ {
		p.wg.Add(1)
		go p.worker(i)
	}
}

// worker 执行任务的 goroutine
func (p *pool) worker(workerID int) {
	defer p.wg.Done()
	for {
		select {
		case task, ok := <-p.taskQueue:
			if !ok {
				// 任务队列已关闭
				return
			}
			// 执行任务
			p.runTask(workerID, task)
		case <-p.ctx.Done():
			return
		}
	}
}

func (p *pool) runTask(workerID int, task Task) {
	defer func() {
		if r := recover(); r != nil {
			// 记录panic信息，可选择记录日志
			atomic.AddInt32(&p.errCount, 1) // 将panic也计入错误

			if p.onErr != nil {
				p.onErr(fmt.Errorf("worker %d panic: %v", workerID, r))
			}

			// 可选：重启worker保持池的worker数量
			if atomic.LoadInt32(&p.closed) == 0 {
				p.wg.Add(1)
				go p.worker(workerID) // 重启一个新的worker替代当前崩溃的worker
			}
		}
	}()
	// 执行任务
	if err := task(p.ctx); err != nil {
		atomic.AddInt32(&p.errCount, 1) // 使用原子操作增加失败任务数
	}
}

// Submit 提交任务到池中
func (p *pool) Submit(task Task) {
	// 如果池已经关闭，直接返回
	if atomic.LoadInt32(&p.closed) == 1 {
		return
	}

	p.taskQueue <- task
}

// Shutdown 关闭 pool，等待所有任务完成并返回失败的任务数
// 将 Cancel 集成到 Shutdown 中，并避免重复关闭
func (p *pool) Shutdown() int32 {
	// 如果池已经关闭，直接返回失败任务数
	if !atomic.CompareAndSwapInt32(&p.closed, 0, 1) {
		return atomic.LoadInt32(&p.errCount)
	}

	close(p.taskQueue) // 关闭任务队列，停止接收新任务
	p.wg.Wait()        // 等待所有任务完成
	p.cancel()         // 取消所有任务的 Context

	return atomic.LoadInt32(&p.errCount)
}
