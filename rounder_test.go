package polygol

import (
	"testing"
)

func TestRounderRound(t *testing.T) {
	var pt1, pt2, pt3 *point

	// no overlap
	rounder.reset()
	pt1 = &point{x: 3, y: 4}
	pt2 = &point{x: 4, y: 5}
	pt3 = &point{x: 5, y: 5}
	expect(t, rounder.round(pt1.x, pt1.y).equal(*pt1))
	expect(t, rounder.round(pt2.x, pt2.y).equal(*pt2))
	expect(t, rounder.round(pt3.x, pt3.y).equal(*pt3))

	// exact overlap
	rounder.reset()
	pt1 = &point{x: 3, y: 4}
	pt2 = &point{x: 4, y: 5}
	pt3 = &point{x: 3, y: 4}
	expect(t, rounder.round(pt1.x, pt1.y).equal(*pt1))
	expect(t, rounder.round(pt2.x, pt2.y).equal(*pt2))
	expect(t, rounder.round(pt3.x, pt3.y).equal(*pt3))

	// rounding one coordinate
	rounder.reset()
	pt1 = &point{x: 3, y: 4}
	pt2 = &point{x: 3 + epsilon, y: 4}
	pt3 = &point{x: 3, y: 4 + epsilon}
	expect(t, rounder.round(pt1.x, pt1.y).equal(*pt1))
	expect(t, rounder.round(pt2.x, pt2.y).equal(*pt1))
	expect(t, rounder.round(pt3.x, pt3.y).equal(*pt1))

	// rounding both coordinates
	rounder.reset()
	pt1 = &point{x: 3, y: 4}
	pt2 = &point{x: 3 + epsilon, y: 4 + epsilon}
	expect(t, rounder.round(pt1.x, pt1.y).equal(*pt1))
	expect(t, rounder.round(pt2.x, pt2.y).equal(*pt1))

	// preseed with 0
	rounder.reset()
	pt1 = &point{x: epsilon / 2, y: -epsilon / 2}
	expect(t, pt1.x != 0)
	expect(t, pt1.y != 0)
	expect(t, rounder.round(pt1.x, pt1.y).equal(point{x: 0, y: 0}))
}
