package deadlock

import "sync"

// ProblemDeadlock demonstrates a potential deadlock when two mutexes are acquired in opposite order.
func ProblemDeadlock() {
	mu1 := &sync.Mutex{}
	mu2 := &sync.Mutex{}
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		mu1.Lock()
		defer mu1.Unlock()
		// Simulate some work
		mu2.Lock()
		defer mu2.Unlock()
	}()
	// Acquire locks in opposite order
	mu2.Lock()
	defer mu2.Unlock()
	mu1.Lock()
	defer mu1.Unlock()
	wg.Wait()
}
