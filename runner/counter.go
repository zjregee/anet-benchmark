package runner

import (
	"fmt"
	"time"
	"sync/atomic"

	"github.com/montanaflynn/stats"
)

type Counter struct {
	total  int64
	failed int64
	costs  []int64
}

func NewCounter() *Counter {
	return &Counter{}
}

func (c *Counter) Idx() int64 {
	return atomic.AddInt64(&c.total, 1) - 1
}

func (c *Counter) Reset(total int64) {
	atomic.StoreInt64(&c.total, 0)
	atomic.StoreInt64(&c.failed, 0)
	c.costs = make([]int64, total)
}

func (c *Counter) AddRecord(idx int64, err error, cost int64) {
	c.costs[idx] = cost
	if err != nil {
		atomic.AddInt64(&c.failed, 1)
	}
}

func (c *Counter) Report(title string, totalns int64, concurrent int, total int64, echoSize int) {
	ms, sec := int64(time.Millisecond), int64(time.Second)
	fmt.Printf("[%s] took %d ms for %d requests\n", title, totalns / ms, c.total)
	fmt.Printf("[%s] requests total: %d, failed: %d\n", title, c.total, c.failed)
	tps := float64(c.total * sec) / float64(totalns)
	costs := make([]float64, len(c.costs))
	for i := range c.costs {
		costs[i] = float64(c.costs[i])
	}
	tp99, _ := stats.Percentile(costs, 99)
	fmt.Printf("[%s]: TPS: %.2f, TP99 %.2f us (c=%d, n=%d, b=%d byte)\n\n", title, tps, tp99, concurrent, total, echoSize)
}
