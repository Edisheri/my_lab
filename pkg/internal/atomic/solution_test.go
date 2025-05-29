package atomic

import "testing"

func TestSolutionAtomicAdd(t *testing.T) {
	n := 1000
	total := SolutionAtomicAdd(n)
	if total != 10000 { // 10 goroutines * n increments each
		t.Errorf("Expected %d, got %d", 10000, total)
	}
}
