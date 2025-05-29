package tasks

import (
	"sync"
	"testing"
	"math/rand"
	"time"
)

func TestSafeQueue(t *testing.T) {
	q := NewSafeQueue()
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func(val int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			q.Push(val)
		}(i)
	}
	wg.Wait()

	popWg := sync.WaitGroup{}
	for i := 0; i < 1000; i++ {
		popWg.Add(1)
		go func() {
			defer popWg.Done()
			_, ok := q.Pop()
			if !ok {
				t.Error("Pop failed")
			}
		}()
	}
	popWg.Wait()
}