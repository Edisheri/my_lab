package cas

import "testing"

func TestSolutionFindMax(t *testing.T) {
	nums := []int64{10, 5, 7, 3, 2, 20, 15}
	max := SolutionFindMax(nums)
	if max != 20 {
		t.Errorf("Expected max 20, got %d", max)
	}
}
