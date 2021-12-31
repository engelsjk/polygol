package polygol

import (
	"math"
	"testing"
)

func TestVectorCrossProduct(t *testing.T) {
	v1 := []float64{1, 2}
	v2 := []float64{3, 4}
	expect(t, crossProduct(v1, v2) == -2.0)
}

func TestVectorDotProduct(t *testing.T) {
	v1 := []float64{1, 2}
	v2 := []float64{3, 4}
	expect(t, dotProduct(v1, v2) == 11.0)
}

func TestVectorLength(t *testing.T) {
	var v []float64

	// horizontal
	v = []float64{3, 0}
	expect(t, length(v) == 3.0)

	// vertical
	v = []float64{0, -2}
	expect(t, length(v) == 2.0)

	// 3-4-5
	v = []float64{3, 4}
	expect(t, length(v) == 5.0)
}

func TestVectorCompareAngles(t *testing.T) {
	var pt1, pt2, pt3 []float64

	// colinear
	pt1 = []float64{1, 1}
	pt2 = []float64{2, 2}
	pt3 = []float64{3, 3}
	expect(t, compareAngles(pt1, pt2, pt3) == 0)
	expect(t, compareAngles(pt2, pt1, pt3) == 0)
	expect(t, compareAngles(pt2, pt3, pt1) == 0)
	expect(t, compareAngles(pt3, pt2, pt1) == 0)

	// offset
	pt1 = []float64{0, 0}
	pt2 = []float64{1, 1}
	pt3 = []float64{1, 0}
	expect(t, compareAngles(pt1, pt2, pt3) == -1)
	expect(t, compareAngles(pt2, pt1, pt3) == 1)
	expect(t, compareAngles(pt2, pt3, pt1) == -1)
	expect(t, compareAngles(pt3, pt2, pt1) == 1)
}

func TestVectorSineAndCosineOfAngle(t *testing.T) {
	var shared, base, angle []float64

	// parallel
	shared = []float64{0, 0}
	base = []float64{1, 0}
	angle = []float64{1, 0}
	expect(t, sineOfAngle(shared, base, angle) == 0.0)
	expect(t, cosineOfAngle(shared, base, angle) == 1.0)

	// 45 degrees
	shared = []float64{0, 0}
	base = []float64{1, 0}
	angle = []float64{1, -1}
	expect(t, almostEqual(sineOfAngle(shared, base, angle), math.Sqrt(2.0)/2.0))
	expect(t, almostEqual(cosineOfAngle(shared, base, angle), math.Sqrt(2.0)/2.0))

	// 90 degrees
	shared = []float64{0, 0}
	base = []float64{1, 0}
	angle = []float64{0, -1}
	expect(t, sineOfAngle(shared, base, angle) == 1)
	expect(t, cosineOfAngle(shared, base, angle) == 0)

	// 135 degrees
	shared = []float64{0, 0}
	base = []float64{1, 0}
	angle = []float64{-1, -1}
	expect(t, almostEqual(sineOfAngle(shared, base, angle), math.Sqrt(2.0)/2.0))
	expect(t, almostEqual(cosineOfAngle(shared, base, angle), -math.Sqrt(2.0)/2.0))

	// anti-parallel
	shared = []float64{0, 0}
	base = []float64{1, 0}
	angle = []float64{-1, 0}
	expect(t, sineOfAngle(shared, base, angle) == 0)
	expect(t, cosineOfAngle(shared, base, angle) == -1)

	// 225 degrees
	shared = []float64{0, 0}
	base = []float64{1, 0}
	angle = []float64{-1, 1}
	expect(t, almostEqual(sineOfAngle(shared, base, angle), -math.Sqrt(2.0)/2.0))
	expect(t, almostEqual(cosineOfAngle(shared, base, angle), -math.Sqrt(2.0)/2.0))

	// 270 degrees
	shared = []float64{0, 0}
	base = []float64{1, 0}
	angle = []float64{0, 1}
	expect(t, sineOfAngle(shared, base, angle) == -1)
	expect(t, cosineOfAngle(shared, base, angle) == 0)

	// 315 degrees
	shared = []float64{0, 0}
	base = []float64{1, 0}
	angle = []float64{1, 1}
	expect(t, almostEqual(sineOfAngle(shared, base, angle), -math.Sqrt(2.0)/2.0))
	expect(t, almostEqual(cosineOfAngle(shared, base, angle), math.Sqrt(2.0)/2.0))
}

func TestVectorPerpindicular(t *testing.T) {
	var v, r []float64

	// vertical
	v = []float64{0, 1}
	r = perpendicular(v)
	expect(t, dotProduct(v, r) == 0)
	expect(t, crossProduct(v, r) != 0)

	// horizontal
	v = []float64{1, 0}
	r = perpendicular(v)
	expect(t, dotProduct(v, r) == 0)
	expect(t, crossProduct(v, r) != 0)

	// 45 degrees
	v = []float64{1, 1}
	r = perpendicular(v)
	expect(t, dotProduct(v, r) == 0)
	expect(t, crossProduct(v, r) != 0)

	// 120 degrees
	v = []float64{-1, 2}
	r = perpendicular(v)
	expect(t, dotProduct(v, r) == 0)
	expect(t, crossProduct(v, r) != 0)
}

func TestVectorClosestPoint(t *testing.T) {
	var pA1, pA2, pB, cp, expected []float64

	// on line
	pA1 = []float64{2, 2}
	pA2 = []float64{3, 3}
	pB = []float64{-1, -1}
	cp = closestPoint(pA1, pA2, pB)
	expect(t, equal(cp, pB))

	// on first point
	pA1 = []float64{2, 2}
	pA2 = []float64{3, 3}
	pB = []float64{2, 2}
	cp = closestPoint(pA1, pA2, pB)
	expect(t, equal(cp, pB))

	// off line above
	pA1 = []float64{2, 2}
	pA2 = []float64{3, 1}
	pB = []float64{3, 7}
	expected = []float64{0, 4}
	expect(t, equal(closestPoint(pA1, pA2, pB), expected))
	expect(t, equal(closestPoint(pA2, pA1, pB), expected))

	// off line below
	pA1 = []float64{2, 2}
	pA2 = []float64{3, 1}
	pB = []float64{0, 2}
	expected = []float64{1, 3}
	expect(t, equal(closestPoint(pA1, pA2, pB), expected))
	expect(t, equal(closestPoint(pA2, pA1, pB), expected))

	// off line perpendicular to first point
	pA1 = []float64{2, 2}
	pA2 = []float64{3, 3}
	pB = []float64{1, 3}
	cp = closestPoint(pA1, pA2, pB)
	expected = []float64{2, 2}
	expect(t, equal(cp, expected))

	// horizontal vector
	pA1 = []float64{2, 2}
	pA2 = []float64{3, 2}
	pB = []float64{1, 3}
	cp = closestPoint(pA1, pA2, pB)
	expected = []float64{1, 2}
	expect(t, equal(cp, expected))

	// vertical vector
	pA1 = []float64{2, 2}
	pA2 = []float64{2, 3}
	pB = []float64{1, 3}
	cp = closestPoint(pA1, pA2, pB)
	expected = []float64{2, 3}
	expect(t, equal(cp, expected))

	// on line but dot product does not think so - part of issue 60-2
	pA1 = []float64{-45.3269382, -1.4059341}
	pA2 = []float64{-45.326737413921656, -1.40635}
	pB = []float64{-45.326833968900424, -1.40615}
	cp = closestPoint(pA1, pA2, pB)
	expect(t, equal(cp, pB))
}

func TestVectorVerticalIntersection(t *testing.T) {
	var pt, i, v []float64
	var x float64

	// horizontal
	pt = []float64{42, 3}
	v = []float64{-2, 0}
	x = 37
	i = verticalIntersection(v, pt, x)
	expect(t, i[0] == 37)
	expect(t, i[1] == 3)

	// vertical
	pt = []float64{42, 3}
	v = []float64{0, 4}
	x = 37
	expect(t, verticalIntersection(v, pt, x) == nil)

	// 45 degree
	pt = []float64{1, 1}
	v = []float64{1, 1}
	x = -2
	i = verticalIntersection(v, pt, x)
	expect(t, i[0] == -2)
	expect(t, i[1] == -2)

	// upper left quadrant
	pt = []float64{-1, 1}
	v = []float64{-2, 1}
	x = -3
	i = verticalIntersection(v, pt, x)
	expect(t, i[0] == -3)
	expect(t, i[1] == 2)
}

func TestVectorHorizontalIntersection(t *testing.T) {
	var pt, i, v []float64
	var y float64

	// horizontal
	pt = []float64{42, 3}
	v = []float64{-2, 0}
	y = 37
	expect(t, horizontalIntersection(v, pt, y) == nil)

	// vertical
	pt = []float64{42, 3}
	v = []float64{0, 4}
	y = 37
	i = horizontalIntersection(v, pt, y)
	expect(t, i[0] == 42)
	expect(t, i[1] == 37)

	// 45 degree
	pt = []float64{1, 1}
	v = []float64{1, 1}
	y = 4
	i = horizontalIntersection(v, pt, y)
	expect(t, i[0] == 4)
	expect(t, i[1] == 4)

	// bottom left quadrant
	pt = []float64{-1, -1}
	v = []float64{-2, -1}
	y = -3
	i = horizontalIntersection(v, pt, y)
	expect(t, i[0] == -5)
	expect(t, i[1] == -3)
}

func TestVectorIntersection(t *testing.T) {
	var i, v1, v2 []float64

	p1 := []float64{42, 42}
	p2 := []float64{-32, 46}

	// parallel
	v1 = []float64{1, 2}
	v2 = []float64{-1, -2}
	i = intersection(v1, v2, p1, p2)
	expect(t, i == nil)

	// horizontal and vertical
	v1 = []float64{0, 2}
	v2 = []float64{-1, 0}
	i = intersection(v1, v2, p1, p2)
	expect(t, i[0] == 42)
	expect(t, i[1] == 46)

	// horizontal
	v1 = []float64{1, 1}
	v2 = []float64{-1, 0}
	i = intersection(v1, v2, p1, p2)
	expect(t, i[0] == 46)
	expect(t, i[1] == 46)

	// vertical
	v1 = []float64{1, 1}
	v2 = []float64{0, 1}
	i = intersection(v1, v2, p1, p2)
	expect(t, i[0] == -32)
	expect(t, i[1] == -32)

	// 45 degree && 135 degree
	v1 = []float64{1, 1}
	v2 = []float64{-1, 1}
	i = intersection(v1, v2, p1, p2)
	expect(t, i[0] == 7)
	expect(t, i[1] == 7)

	// consistency
	// Taken from https://github.com/mfogel/polygon-clipping/issues/37
	p1 = []float64{0.523787, 51.281453}
	v1 = []float64{0.0002729999999999677, 0.0002729999999999677}
	p2 = []float64{0.523985, 51.281651}
	v2 = []float64{0.000024999999999941735, 0.000049000000004184585}
	i1 := intersection(v1, v2, p1, p2)
	i2 := intersection(v2, v1, p2, p1)
	expect(t, i1[0] == i2[0])
	expect(t, i1[1] == i2[1])
}
