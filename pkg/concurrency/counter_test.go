package concurrency

import "testing"

func TestCounter(t *testing.T) {
	c := Counter{}
	for i := 0; i < 100; i++ {
		c.Increment()
	}
	if c.Value() != 100 {
		t.Errorf("Expected 100, got %d", c.Value())
	}
}
