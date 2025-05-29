package barrier

import (
	"testing"
)

// TestSolutionBarrier just ensures SolutionBarrier completes without deadlock.
func TestSolutionBarrier(t *testing.T) {
	// Running the barrier should not deadlock.
	SolutionBarrier(3)
}
