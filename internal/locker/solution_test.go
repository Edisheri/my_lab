package locker

import "testing"

func TestSolutionLockerIncrement(t *testing.T) {
	n := 1000
	total := SolutionLockerIncrement(n)
	if total != 10000 {
		t.Errorf("Expected %d, got %d", 10000, total)
	}
}
