package polygol

import "math"

const (
	epsilon = float64(7.)/3 - float64(4.)/3 - float64(1.)
	// epsilon = 2e-12
)

var (
	epsilonSq = epsilon * epsilon
)

func flpCmp(a, b float64) int {
	// check if they're both 0
	if -epsilon < a && a < epsilon {
		if -epsilon < b && b < epsilon {
			return 0
		}
	}

	// check if they're flp equal
	ab := a - b
	if ab*ab < epsilonSq*a*b {
		return 0
	}

	if a < b {
		return -1
	}

	return 1
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) <= epsilon
}
