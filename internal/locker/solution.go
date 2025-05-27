package locker

import "sync"

// MyLocker is a simple implementation of sync.Locker using a Mutex.
type MyLocker struct {
	mu sync.Mutex
}

// Lock locks the underlying mutex.
func (l *MyLocker) Lock() {
	l.mu.Lock()
}

// Unlock unlocks the underlying mutex.
func (l *MyLocker) Unlock() {
	l.mu.Unlock()
}

// SolutionLockerIncrement uses a custom Locker to increment a counter.
func SolutionLockerIncrement(n int) int {
	var counter int
	var locker MyLocker
	var wg sync.WaitGroup
	// Launch 10 goroutines to increment the counter n times each using the custom locker.
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < n; j++ {
				locker.Lock()
				counter++
				locker.Unlock()
			}
		}()
	}
	wg.Wait()
	return counter
}
