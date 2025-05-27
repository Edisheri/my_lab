package concurrency

import "sync/atomic"

// Counter provides a thread-safe counter using atomic operations.
type Counter struct {
	value int64
}

// Increment increases the counter by 1.
func (c *Counter) Increment() {
	atomic.AddInt64(&c.value, 1)
}

// Value returns the current value of the counter.
func (c *Counter) Value() int64 {
	return atomic.LoadInt64(&c.value)
}
