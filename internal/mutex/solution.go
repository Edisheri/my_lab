package mutex

import "sync"

// SolutionMutexIncrement increments a counter from multiple goroutines using sync.Mutex to avoid data races.
func SolutionMutexIncrement(n int) int {
	var counter int
	var mu sync.Mutex
	var wg sync.WaitGroup
	// Launch 10 goroutines to increment the counter n times each.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < n; j++ {
				mu.Lock()
				counter++
				mu.Unlock()
			}
		}()
	}
	wg.Wait()
	return counter
}
