package mutex

import "testing"

func TestSolutionMutexIncrement(t *testing.T) {
	n := 1000
	total := SolutionMutexIncrement(n)
	if total != 10000 {
		t.Errorf("Expected %d, got %d", 10000, total)
	}
}
