package barrier

import (
	"fmt"
	"sync"
)

// SolutionBarrier launches n goroutines that wait for a signal before proceeding.
func SolutionBarrier(n int) {
	var wg sync.WaitGroup
	start := make(chan struct{})
	wg.Add(n)
	for i := 0; i < n; i++ {
		go func(id int) {
			defer wg.Done()
			<-start // wait for signal to start
			fmt.Printf("Goroutine %d running\n", id)
		}(i)
	}
	// Close the channel to signal all goroutines at once
	close(start)
	wg.Wait()
}
