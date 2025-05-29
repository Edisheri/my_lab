package tasks

import (
	"sync"
	"testing"
	"math/rand"
	"time"
)

func TestUniqueUsers(t *testing.T) {
	uu := NewUniqueUsers()
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		id := rand.Intn(100)
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(5)) * time.Millisecond)
			uu.AddUser(id)
		}(id)
	}
	wg.Wait()
	count := uu.Count()
	if count < 1 || count > 100 {
		t.Errorf("Expected count in [1,100], got %d", count)
	}
}