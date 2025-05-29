package deadlock

import "sync"

// SolutionDeadlock avoids the deadlock by locking mutexes in a consistent order.
func SolutionDeadlock() {
	mu1 := &sync.Mutex{}
	mu2 := &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		// Lock in consistent order: mu1, then mu2
		mu1.Lock()
		mu2.Lock()
		mu2.Unlock()
		mu1.Unlock()
	}()
	// Lock in the same order: mu1, then mu2
	mu1.Lock()
	mu2.Lock()
	mu2.Unlock()
	mu1.Unlock()
	wg.Wait()
}
