package polygol

import "testing"

func TestFlpCompare(t *testing.T) {
	var a, b float64

	// exactly equal
	a = 1
	b = 1
	expect(t, flpCmp(a, b) == 0)

	// flp equal
	a = 1
	b = 1 + epsilon
	expect(t, flpCmp(a, b) == 0)

	// barely less than
	a = 1
	b = 1 + epsilon*2
	expect(t, flpCmp(a, b) == -1)

	// less than
	a = 1
	b = 2
	expect(t, flpCmp(a, b) == -1)

	// barely more than
	a = 1 + epsilon*2
	b = 1
	expect(t, flpCmp(a, b) == 1)

	// more than
	a = 2
	b = 1
	expect(t, flpCmp(a, b) == 1)

	// both flp equal to 0
	a = 0.0
	b = epsilon - epsilon*epsilon
	expect(t, flpCmp(a, b) == 0)

	// really close to 0
	a = epsilon
	b = epsilon + epsilon*epsilon*2
	expect(t, flpCmp(a, b) == -1)
}
