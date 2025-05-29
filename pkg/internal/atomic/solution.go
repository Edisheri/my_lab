package atomic

import (
	"sync"
	"sync/atomic"
)

// SolutionAtomicAdd increments a counter from multiple goroutines using atomic operations to avoid data races.
func SolutionAtomicAdd(n int) int {
	var counter int64
	var wg sync.WaitGroup
	// Launch 10 goroutines to increment the counter n times each.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < n; j++ {
				atomic.AddInt64(&counter, 1)
			}
		}()
	}
	wg.Wait()
	return int(counter)
}
