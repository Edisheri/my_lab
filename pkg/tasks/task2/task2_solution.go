package tasks

import "sync"

type SafeQueue struct {
	mu    sync.Mutex
	queue []int
}

func NewSafeQueue() *SafeQueue {
	return &SafeQueue{queue: make([]int, 0)}
}

func (q *SafeQueue) Push(val int) {
	q.mu.Lock()
	q.queue = append(q.queue, val)
	q.mu.Unlock()
}

func (q *SafeQueue) Pop() (int, bool) {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.queue) == 0 {
		return 0, false
	}
	val := q.queue[0]
	q.queue = q.queue[1:]
	return val, true
}