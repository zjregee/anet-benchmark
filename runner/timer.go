package runner

import (
	"sync"
	"time"
	"sync/atomic"
)

func NewTimer(window time.Duration) *Timer {
	t := &Timer{ window: window }
	t.refresh()
	return t
}

type Timer struct {
	sync.Once
	now int64
	window time.Duration
}

func (t *Timer) Now() int64 {
	return atomic.LoadInt64(&t.now)
}

func (t *Timer) refresh() {
	t.Do(func() {
		atomic.StoreInt64(&t.now, time.Now().UnixNano())
		go func() {
			for now := range time.Tick(t.window) {
				atomic.StoreInt64(&t.now, now.UnixNano())
			}
		}()
	})
}
