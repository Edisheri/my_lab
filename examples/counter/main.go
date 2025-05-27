package main

import (
	"fmt"
	"my_lab/pkg/concurrency"
)

func main() {
	c := concurrency.Counter{}
	// Increment counter
	c.Increment()
	c.Increment()
	fmt.Printf("Counter value: %d\n", c.Value())
}
