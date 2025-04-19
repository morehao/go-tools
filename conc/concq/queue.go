package concq

import (
	"context"
	"errors"
	"sync"
)

// Task 表示一个待处理的任务
type Task func(ctx context.Context) error

// Queue 是一个基于生产者-消费者模型的并发控制器
type Queue struct {
	taskCh  chan Task
	wg      sync.WaitGroup
	ctx     context.Context
	cancel  context.CancelFunc
	workerN int
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

// Submit 提交一个任务到队列
// 如果队列已关闭，会返回错误
func (q *Queue) Submit(t Task) error {
	select {
	case q.taskCh <- t:
		return nil
	case <-q.ctx.Done():
		return errors.New("queue has been shutdown")
	}
}

// start 启动所有 worker
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
			_ = task(q.ctx) // 忽略错误，可根据需求扩展
		}
	}
}

// Shutdown 优雅关闭队列，等待所有任务完成
func (q *Queue) Shutdown() {
	q.cancel()
	close(q.taskCh)
	q.wg.Wait()
}
