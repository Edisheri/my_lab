package main

import (
	"fmt"
	"my_lab/internal/atomic"
	"my_lab/internal/barrier"
	"my_lab/internal/cas"
	"my_lab/internal/deadlock"
	"my_lab/internal/locker"
	"my_lab/internal/mutex"
	"my_lab/internal/queue"
)

func main() {
	fmt.Println("Running concurrency tasks demonstration:")
	// Atomic increment
	fmt.Printf("Atomic Add: %d\n", atomic.SolutionAtomicAdd(1000))
	// CAS max
	nums := []int64{10, 5, 7, 3, 2, 20, 15}
	fmt.Printf("Atomic CAS max: %d\n", cas.SolutionFindMax(nums))
	// Mutex increment
	fmt.Printf("Mutex Add: %d\n", mutex.SolutionMutexIncrement(1000))
	// Locker pattern
	fmt.Printf("Locker Add: %d\n", locker.SolutionLockerIncrement(1000))
	// Queue usage
	q := queue.NewSafeQueue()
	for i := 1; i <= 5; i++ {
		q.Enqueue(i)
	}
	fmt.Printf("Queue Length after enqueue: %d\n", q.Len())
	// Barrier example
	barrier.SolutionBarrier(3)
	// Deadlock solution (no output expected, just ensure no deadlock)
	deadlock.SolutionDeadlock()
	fmt.Println("Deadlock solution executed without deadlock.")
}
