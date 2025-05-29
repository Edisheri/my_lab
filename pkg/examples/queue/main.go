package main

import (
	"fmt"
	"my_lab/internal/queue"
)

func main() {
	q := queue.NewSafeQueue()
	q.Enqueue(10)
	q.Enqueue(20)
	fmt.Println("Dequeued:", q.Dequeue())
	fmt.Println("Queue length:", q.Len())
}
