package concpool

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"
)

func TestPool(t *testing.T) {
	t.Run("normal tasks", func(t *testing.T) {
		var counter int32
		p := New(3, 10)

		for i := 0; i < 5; i++ {
			p.Submit(func(ctx context.Context) error {
				atomic.AddInt32(&counter, 1)
				return nil
			})
		}

		errCount := p.StopAndWait()
		if errCount != 0 {
			t.Errorf("expected 0 errors, got %d", errCount)
		}
		if counter != 5 {
			t.Errorf("expected counter to be 5, got %d", counter)
		}
	})

	t.Run("tasks with errors", func(t *testing.T) {
		var success int32
		p := New(2, 10)

		for i := 0; i < 3; i++ {
			p.Submit(func(ctx context.Context) error {
				atomic.AddInt32(&success, 1)
				return nil
			})
		}
		for i := 0; i < 2; i++ {
			p.Submit(func(ctx context.Context) error {
				return errors.New("task error")
			})
		}

		errCount := p.StopAndWait()
		if errCount != 2 {
			t.Errorf("expected 2 errors, got %d", errCount)
		}
		if success != 3 {
			t.Errorf("expected 3 successes, got %d", success)
		}
	})

	t.Run("panic in task", func(t *testing.T) {
		p := New(1, 5)
		p.Submit(func(ctx context.Context) error {
			panic("boom")
		})
		p.Submit(func(ctx context.Context) error {
			return nil
		})

		errCount := p.StopAndWait()
		if errCount != 1 {
			t.Errorf("expected 1 error due to panic, got %d", errCount)
		}
	})

	t.Run("submit after shutdown", func(t *testing.T) {
		var counter int32
		p := New(1, 2)
		p.Submit(func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		})
		_ = p.StopAndWait()

		// 此时应忽略任务
		p.Submit(func(ctx context.Context) error {
			atomic.AddInt32(&counter, 1)
			return nil
		})
		time.Sleep(100 * time.Millisecond)

		if counter != 1 {
			t.Errorf("expected 1 execution, got %d", counter)
		}
	})

	t.Run("StopAndWait is idempotent", func(t *testing.T) {
		p := New(1, 1)
		p.Submit(func(ctx context.Context) error {
			return nil
		})
		err1 := p.StopAndWait()
		err2 := p.StopAndWait()
		if err1 != err2 {
			t.Errorf("expected same error count on repeated StopAndWait calls, got %d and %d", err1, err2)
		}
	})
}
