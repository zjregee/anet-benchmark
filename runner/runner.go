package runner

import (
	"sync"
	"time"
	"errors"
	"context"
)

type Runner struct {
	timer   *Timer
	counter *Counter
}

func NewRunner() *Runner {
	return &Runner{
		counter: NewCounter(),
		timer: NewTimer(time.Millisecond),
	}
}

func (r *Runner) Run(title string, once EchoOnce, concurrent int, total int64, echoSize int) {
	start := r.timer.Now()
	r.benching(once, concurrent, total, echoSize)
	end := r.timer.Now()
	r.counter.Report(title, end - start, concurrent, total, echoSize)
}

func (r *Runner) echoWithTimeout(ctx context.Context, once EchoOnce, req *Message) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	done := make(chan error, 1)
	go func() {
		_, err := once(req)
		done <- err
	}()
	select {
	case err := <-done:
		return err
	case <-ctx.Done():
		return errors.New("echo timeout")
	}
}

func (r *Runner) benching(once EchoOnce, concurrent int, total int64, echoSize int) {
	var wg sync.WaitGroup
	wg.Add(concurrent)
	r.counter.Reset(total)
	for i := 0; i < concurrent; i++ {
		go func() {
			defer wg.Done()
			ctx := context.Background()
			body := make([]byte, echoSize)
			req := &Message{Message: string(body)}
			for {
				idx := r.counter.Idx()
				if idx >= total {
					return
				}
				begin := r.timer.Now()
				err := r.echoWithTimeout(ctx, once, req)
				end := r.timer.Now()
				cost := end - begin
				r.counter.AddRecord(idx, err, cost)
			}
		}()
	}
	wg.Wait()
	r.counter.total = total
}
