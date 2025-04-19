package concqueue

import (
	"context"
	"errors"
	"sync"
	"sync/atomic"
	"time"
)

type Queue interface {
	Submit(t Task)
	Shutdown() int32
}

// Task 表示一个可执行的任务
type Task func(ctx context.Context) error

// queue 是一个基于生产者-消费者模型的并发控制器
type queue struct {
	taskCh          chan Task
	wg              sync.WaitGroup
	ctx             context.Context
	cancel          context.CancelFunc
	workerN         int
	errCount        int32
	closed          int32
	submitTimeout   time.Duration
	shutdownTimeout time.Duration
	errorHandler    ErrorHandler
	panicHandler    PanicHandler
	logger          Logger
}

// New 创建一个新的 queue 实例
func New(workerCount, queueSize int, options ...Option) Queue {
	ctx, cancel := context.WithCancel(context.Background())
	q := &queue{
		taskCh:          make(chan Task, queueSize),
		ctx:             ctx,
		cancel:          cancel,
		workerN:         workerCount,
		logger:          &defaultLogger{}, // 使用默认 Logger
		submitTimeout:   0,                // 默认不超时
		shutdownTimeout: 0,                // 默认不超时
	}
	for _, opt := range options {
		opt(q)
	}
	q.start()
	return q
}

// Submit (生产者)提交一个任务到队列
func (q *queue) Submit(t Task) {
	if atomic.LoadInt32(&q.closed) == 1 {
		// 队列已关闭，直接丢弃任务
		return
	}

	// 如果设置了超时时间，则使用超时机制
	if q.submitTimeout > 0 {
		timer := time.NewTimer(q.submitTimeout)
		defer timer.Stop()
		select {
		case q.taskCh <- t:
			// 任务提交成功
			return
		case <-q.ctx.Done():
			// 队列已关闭，丢弃任务
			return
		case <-timer.C:
			// 提交超时
			if q.errorHandler != nil {
				q.errorHandler(errors.New("submit timeout"))
			}
			return
		}
	}

	// 没有设置超时时间，则直接提交任务
	select {
	case q.taskCh <- t:
		// 任务提交成功
	case <-q.ctx.Done():
		// 队列已关闭，丢弃任务
		return
	}
}

// start 启动 worker 协程（消费者）
func (q *queue) start() {
	for i := 0; i < q.workerN; i++ {
		q.wg.Add(1)
		go q.worker()
	}
}

// worker 是消费任务的协程
func (q *queue) worker() {
	defer q.wg.Done()
	defer func() {
		if r := recover(); r != nil {
			// 捕获 panic
			if q.panicHandler != nil {
				q.panicHandler(r)
			} else {
				q.logger.Printf("panic occurred: %v", r) // 默认记录 panic 日志
			}
		}
	}()
	for {
		select {
		case <-q.ctx.Done():
			return
		case task, ok := <-q.taskCh:
			if !ok {
				return
			}
			if err := task(q.ctx); err != nil {
				atomic.AddInt32(&q.errCount, 1)
				if q.errorHandler != nil {
					q.errorHandler(err) // 调用 ErrorHandler
				}
			}
		}
	}
}

// Shutdown 主动终止队列，不再接受新任务，并等待所有 worker 停止
func (q *queue) Shutdown() int32 {
	if !atomic.CompareAndSwapInt32(&q.closed, 0, 1) {
		return atomic.LoadInt32(&q.errCount) // 已关闭
	}

	close(q.taskCh)

	// 如果设置了超时时间，则使用超时机制
	if q.shutdownTimeout > 0 {
		done := make(chan struct{})
		go func() {
			defer close(done)
			q.wg.Wait()
		}()
		select {
		case <-done:
			// 所有 worker 都已完成
		case <-time.After(q.shutdownTimeout):
			// 超时
			q.logger.Printf("shutdown timeout")
		}
	} else {
		// 没有设置超时时间，则直接等待所有 worker 完成
		q.wg.Wait()
	}

	q.cancel()
	return atomic.LoadInt32(&q.errCount)
}
