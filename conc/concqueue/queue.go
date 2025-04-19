package concqueue

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
)

// Task 表示一个可执行的任务
type Task func(ctx context.Context) error

// Queue 是一个基于生产者-消费者模型的并发控制器
type Queue struct {
	taskCh   chan Task
	wg       sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	workerN  int
	errCount int64
	closed   int32
}

// New 创建一个新的 Queue 实例
func New(workerCount, queueSize int) *Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &Queue{
		taskCh:  make(chan Task, queueSize),
		ctx:     ctx,
		cancel:  cancel,
		workerN: workerCount,
	}
	q.start()
	return q
}

// Submit (生产者)提交一个任务到队列
// 如果队列已关闭，会返回错误
func (q *Queue) Submit(t Task) error {
	if atomic.LoadInt32(&q.closed) == 1 {
		return errors.New("queue has been shutdown")
	}
	select {
	case q.taskCh <- t:
		return nil
	case <-q.ctx.Done():
		return errors.New("queue has been shutdown")
	}
}

// start 启动 worker 协程（消费者）
func (q *Queue) start() {
	for i := 0; i < q.workerN; i++ {
		q.wg.Add(1)
		go q.worker()
	}
}

// worker 是消费任务的协程
func (q *Queue) worker() {
	defer q.wg.Done()
	for {
		select {
		case <-q.ctx.Done():
			return
		case task, ok := <-q.taskCh:
			if !ok {
				return
			}
			if err := task(q.ctx); err != nil {
				atomic.AddInt64(&q.errCount, 1)
			}
		}
	}
}

// Shutdown 主动终止队列，不再接受新任务，并等待所有 worker 停止
func (q *Queue) Shutdown() int {
	if !atomic.CompareAndSwapInt32(&q.closed, 0, 1) {
		return int(q.errCount) // 已关闭
	}

	close(q.taskCh)
	// 等待 worker 处理完所有任务后退出
	q.wg.Wait()
	q.cancel()
	return int(q.errCount)
}
