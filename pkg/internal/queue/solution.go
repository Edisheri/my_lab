package queue

import "sync"

// SafeQueue is a thread-safe FIFO queue for integers.
type SafeQueue struct {
	items []int
	mu    sync.Mutex
}

// NewSafeQueue creates a new SafeQueue.
func NewSafeQueue() *SafeQueue {
	return &SafeQueue{
		items: make([]int, 0),
	}
}

// Enqueue adds an item to the queue.
func (q *SafeQueue) Enqueue(item int) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.items = append(q.items, item)
}

// Dequeue removes and returns an item from the queue.
// If the queue is empty, Dequeue returns -1.
func (q *SafeQueue) Dequeue() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	if len(q.items) == 0 {
		return -1
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item
}

// Len returns the number of items in the queue.
func (q *SafeQueue) Len() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.items)
}
