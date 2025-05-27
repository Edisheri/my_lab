package cas

import (
	"sync"
	"sync/atomic"
)

// SolutionFindMax finds the maximum value in nums using multiple goroutines and atomic.CompareAndSwap.
func SolutionFindMax(nums []int64) int64 {
	var max int64
	var wg sync.WaitGroup
	for _, n := range nums {
		wg.Add(1)
		go func(val int64) {
			defer wg.Done()
			for {
				old := atomic.LoadInt64(&max)
				if val <= old {
					break
				}
				if atomic.CompareAndSwapInt64(&max, old, val) {
					break
				}
			}
		}(n)
	}
	wg.Wait()
	return max
}
