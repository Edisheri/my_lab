package queue

import (
	"sync"
	"testing"
)

func TestSafeQueue(t *testing.T) {
	q := NewSafeQueue()
	// Enqueue items concurrently
	var wg sync.WaitGroup
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			q.Enqueue(val)
		}(i)
	}
	wg.Wait()
	// Now check the length
	if q.Len() != 5 {
		t.Errorf("Expected queue length 5, got %d", q.Len())
	}
}
