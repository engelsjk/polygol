package polygol

import (
	"math"
	"testing"
)

func equalMultiPoly(m1, m2 [][][][]float64) bool {
	if len(m1) != len(m2) {
		return false
	}
	for i, p := range m1 {
		if len(p) != len(m2[i]) {
			return false
		}
		if !equalPoly(p, m2[i]) {
			return false
		}
	}
	return true
}

func equalPoly(p1, p2 [][][]float64) bool {
	if len(p1) != len(p2) {
		return false
	}
	for i, r := range p1 {
		if len(r) != len(p2[i]) {
			return false
		}
		if !equalRing(r, p2[i]) {
			return false
		}
	}
	return true
}

func equalRing(r1, r2 [][]float64) bool {
	if len(r1) != len(r2) {
		return false
	}
	for i, pt := range r1 {
		if len(pt) != len(r2[i]) {
			return false
		}
		for j := range pt {
			if math.Abs(pt[j]-r2[i][j]) > epsilon {
				return false
			}
		}
	}
	return true
}

// ring
func TestGeomOutRingSimpleTriangle(t *testing.T) {
	t.Parallel()

	// simple triangle

	op := newOperation("")

	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)
	p3 := newPoint(0, 1)
	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p1, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	rings, err := newRingOutFromSegments([]*segment{seg1, seg2, seg3})
	terr(t, err)

	expect(t, len(rings) == 1)
	expect(t, equalRing(rings[0].getGeom(), [][]float64{{0, 0}, {1, 1}, {0, 1}, {0, 0}}))
}

func TestGeomOutRingBowTie(t *testing.T) {
	t.Parallel()

	// bow tie

	op := newOperation("")

	p1 := newPoint(0, 0)
	p2 := newPoint(1, 1)
	p3 := newPoint(0, 2)

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p1, &ringIn{})
	terr(t, err)

	p4 := newPoint(2, 0)
	p5 := p2
	p6 := newPoint(2, 2)

	seg4, err := op.newSegmentFromRing(p4, p5, &ringIn{})
	terr(t, err)
	seg5, err := op.newSegmentFromRing(p5, p6, &ringIn{})
	terr(t, err)
	seg6, err := op.newSegmentFromRing(p6, p4, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	seg4.forceInResult, seg4.inResult = true, true
	seg5.forceInResult, seg5.inResult = true, true
	seg6.forceInResult, seg6.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{
		seg1, seg2, seg3, seg4,
		seg5, seg6,
	})
	terr(t, err)

	expect(t, len(rings) == 2)
	expect(t, equalRing(rings[0].getGeom(), [][]float64{{0, 0}, {1, 1}, {0, 2}, {0, 0}}))
	expect(t, equalRing(rings[1].getGeom(), [][]float64{{1, 1}, {2, 0}, {2, 2}, {1, 1}}))
}

func TestGeomOutRingRinged(t *testing.T) {
	t.Parallel()

	// ring ringed

	op := newOperation("")

	p1 := newPoint(0, 0)
	p2 := newPoint(3, -3)
	p3 := newPoint(3, 0)
	p4 := newPoint(3, 3)

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p4, &ringIn{})
	terr(t, err)
	seg4, err := op.newSegmentFromRing(p4, p1, &ringIn{})
	terr(t, err)

	p5 := newPoint(2, -1)
	p6 := p3
	p7 := newPoint(2, 1)

	seg5, err := op.newSegmentFromRing(p5, p6, &ringIn{})
	terr(t, err)
	seg6, err := op.newSegmentFromRing(p6, p7, &ringIn{})
	terr(t, err)
	seg7, err := op.newSegmentFromRing(p7, p5, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	seg4.forceInResult, seg4.inResult = true, true
	seg5.forceInResult, seg5.inResult = true, true
	seg6.forceInResult, seg6.inResult = true, true
	seg7.forceInResult, seg7.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{
		seg1, seg2, seg3, seg4,
		seg5, seg6, seg7,
	})
	terr(t, err)

	expect(t, len(rings) == 2)
	expect(t, equalRing(rings[0].getGeom(), [][]float64{{3, 0}, {2, 1}, {2, -1}, {3, 0}}))
	expect(t, equalRing(rings[1].getGeom(), [][]float64{{0, 0}, {3, -3}, {3, 3}, {0, 0}}))
}

func TestGeomOutRingRingedInterior(t *testing.T) {
	t.Parallel()

	// ringed ring interior ring starting point extraneous

	op := newOperation("")

	p1 := &point{x: 0, y: 0}
	p2 := &point{x: 5, y: -5}
	p3 := &point{x: 4, y: 0}
	p4 := &point{x: 5, y: 5}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p4, &ringIn{})
	terr(t, err)
	seg4, err := op.newSegmentFromRing(p4, p1, &ringIn{})
	terr(t, err)

	p5 := &point{x: 1, y: 0}
	p6 := &point{x: 4, y: 1}
	p7 := p3
	p8 := &point{x: 4, y: -1}

	seg5, err := op.newSegmentFromRing(p5, p6, &ringIn{})
	terr(t, err)
	seg6, err := op.newSegmentFromRing(p6, p7, &ringIn{})
	terr(t, err)
	seg7, err := op.newSegmentFromRing(p7, p8, &ringIn{})
	terr(t, err)
	seg8, err := op.newSegmentFromRing(p8, p5, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	seg4.forceInResult, seg4.inResult = true, true
	seg5.forceInResult, seg5.inResult = true, true
	seg6.forceInResult, seg6.inResult = true, true
	seg7.forceInResult, seg7.inResult = true, true
	seg8.forceInResult, seg8.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{
		seg1, seg2, seg3, seg4,
		seg5, seg6, seg7, seg8,
	})
	terr(t, err)

	expect(t, len(rings) == 2)
	expect(t, equalRing(rings[0].getGeom(), [][]float64{{4, 1}, {1, 0}, {4, -1}, {4, 1}}))
	expect(t, equalRing(rings[1].getGeom(), [][]float64{{0, 0}, {5, -5}, {4, 0}, {5, 5}, {0, 0}}))
}

func TestGeomOutRingRingedBowTie(t *testing.T) {
	t.Parallel()

	// ringed ring and bow tie at same point

	op := newOperation("")

	p1 := &point{x: 0, y: 0}
	p2 := &point{x: 3, y: -3}
	p3 := &point{x: 3, y: 0}
	p4 := &point{x: 3, y: 3}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p4, &ringIn{})
	terr(t, err)
	seg4, err := op.newSegmentFromRing(p4, p1, &ringIn{})
	terr(t, err)

	p5 := &point{x: 2, y: -1}
	p6 := p3
	p7 := &point{x: 2, y: 1}

	seg5, err := op.newSegmentFromRing(p5, p6, &ringIn{})
	terr(t, err)
	seg6, err := op.newSegmentFromRing(p6, p7, &ringIn{})
	terr(t, err)
	seg7, err := op.newSegmentFromRing(p7, p5, &ringIn{})
	terr(t, err)

	p8 := p3
	p9 := &point{x: 4, y: -1}
	p10 := &point{x: 4, y: 1}

	seg8, err := op.newSegmentFromRing(p8, p9, &ringIn{})
	terr(t, err)
	seg9, err := op.newSegmentFromRing(p9, p10, &ringIn{})
	terr(t, err)
	seg10, err := op.newSegmentFromRing(p10, p8, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	seg4.forceInResult, seg4.inResult = true, true
	seg5.forceInResult, seg5.inResult = true, true
	seg6.forceInResult, seg6.inResult = true, true
	seg7.forceInResult, seg7.inResult = true, true
	seg8.forceInResult, seg8.inResult = true, true
	seg9.forceInResult, seg9.inResult = true, true
	seg10.forceInResult, seg10.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{
		seg1, seg2, seg3, seg4, seg5,
		seg6, seg7, seg8, seg9, seg10,
	})
	terr(t, err)

	expect(t, len(rings) == 3)
	expect(t, equalRing(rings[0].getGeom(), [][]float64{{3, 0}, {2, 1}, {2, -1}, {3, 0}}))
	expect(t, equalRing(rings[1].getGeom(), [][]float64{{0, 0}, {3, -3}, {3, 3}, {0, 0}}))
	expect(t, equalRing(rings[2].getGeom(), [][]float64{{3, 0}, {4, -1}, {4, 1}, {3, 0}}))
}

func TestGeomOutRingDoubleBowTie(t *testing.T) {
	t.Parallel()

	// double bow tie

	op := newOperation("")

	p1 := &point{x: 0, y: 0}
	p2 := &point{x: 1, y: -2}
	p3 := &point{x: 1, y: 2}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p1, &ringIn{})
	terr(t, err)

	p4 := p2
	p5 := &point{x: 2, y: -3}
	p6 := &point{x: 2, y: -1}

	seg4, err := op.newSegmentFromRing(p4, p5, &ringIn{})
	terr(t, err)
	seg5, err := op.newSegmentFromRing(p5, p6, &ringIn{})
	terr(t, err)
	seg6, err := op.newSegmentFromRing(p6, p4, &ringIn{})
	terr(t, err)

	p7 := p3
	p8 := &point{x: 2, y: 1}
	p9 := &point{x: 2, y: 3}

	seg7, err := op.newSegmentFromRing(p7, p8, &ringIn{})
	terr(t, err)
	seg8, err := op.newSegmentFromRing(p8, p9, &ringIn{})
	terr(t, err)
	seg9, err := op.newSegmentFromRing(p9, p7, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	seg4.forceInResult, seg4.inResult = true, true
	seg5.forceInResult, seg5.inResult = true, true
	seg6.forceInResult, seg6.inResult = true, true
	seg7.forceInResult, seg7.inResult = true, true
	seg8.forceInResult, seg8.inResult = true, true
	seg9.forceInResult, seg9.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{
		seg1, seg2, seg3, seg4,
		seg5, seg6, seg7, seg8, seg9,
	})
	terr(t, err)

	expect(t, len(rings) == 3)
	expect(t, equalRing(rings[0].getGeom(), [][]float64{{0, 0}, {1, -2}, {1, 2}, {0, 0}}))
	expect(t, equalRing(rings[1].getGeom(), [][]float64{{1, -2}, {2, -3}, {2, -1}, {1, -2}}))
	expect(t, equalRing(rings[2].getGeom(), [][]float64{{1, 2}, {2, 1}, {2, 3}, {1, 2}}))
}

func TestGeomOutRingDoubleRinged(t *testing.T) {
	t.Parallel()

	// double ringed ring

	op := newOperation("")

	p1 := &point{x: 0, y: 0}
	p2 := &point{x: 5, y: -5}
	p3 := &point{x: 5, y: 5}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p1, &ringIn{})
	terr(t, err)

	p4 := &point{x: 1, y: -1}
	p5 := p2
	p6 := &point{x: 2, y: -1}

	seg4, err := op.newSegmentFromRing(p4, p5, &ringIn{})
	terr(t, err)
	seg5, err := op.newSegmentFromRing(p5, p6, &ringIn{})
	terr(t, err)
	seg6, err := op.newSegmentFromRing(p6, p4, &ringIn{})
	terr(t, err)

	p7 := &point{x: 1, y: 1}
	p8 := p3
	p9 := &point{x: 2, y: 1}

	seg7, err := op.newSegmentFromRing(p7, p8, &ringIn{})
	terr(t, err)
	seg8, err := op.newSegmentFromRing(p8, p9, &ringIn{})
	terr(t, err)
	seg9, err := op.newSegmentFromRing(p9, p7, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	seg4.forceInResult, seg4.inResult = true, true
	seg5.forceInResult, seg5.inResult = true, true
	seg6.forceInResult, seg6.inResult = true, true
	seg7.forceInResult, seg7.inResult = true, true
	seg8.forceInResult, seg8.inResult = true, true
	seg9.forceInResult, seg9.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{
		seg1, seg2, seg3, seg4,
		seg5, seg6, seg7, seg8, seg9,
	})
	terr(t, err)

	expect(t, len(rings) == 3)
	expect(t, equalRing(rings[0].getGeom(), [][]float64{{5, -5}, {2, -1}, {1, -1}, {5, -5}}))
	expect(t, equalRing(rings[1].getGeom(), [][]float64{{5, 5}, {1, 1}, {2, 1}, {5, 5}}))
	expect(t, equalRing(rings[2].getGeom(), [][]float64{{0, 0}, {5, -5}, {5, 5}, {0, 0}}))
}

func TestGeomOutRingMalformed(t *testing.T) {
	t.Parallel()

	// errors on on malformed ring

	op := newOperation("")

	p1 := &point{x: 0, y: 0}
	p2 := &point{x: 1, y: 1}
	p3 := &point{x: 0, y: 1}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p1, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, false

	_, err = newRingOutFromSegments([]*segment{seg1, seg2, seg3})
	expect(t, err != nil)
}

func TestGeomOutRingExterior(t *testing.T) {
	t.Parallel()

	// exterior ring

	op := newOperation("")

	p1 := &point{x: 0, y: 0}
	p2 := &point{x: 1, y: 1}
	p3 := &point{x: 0, y: 1}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p1, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{seg1, seg2, seg3})
	terr(t, err)

	ring := rings[0]

	expect(t, ring.calcEnclosingRing() == nil)
	expect(t, ring.calcIsExteriorRing())
	expect(t, equalRing(ring.getGeom(), [][]float64{{0, 0}, {1, 1}, {0, 1}, {0, 0}}))
}

func TestGeomOutRingInteriorReverse(t *testing.T) {
	t.Parallel()

	// interior ring points reversed

	op := newOperation("")

	p1 := &point{x: 0, y: 0}
	p2 := &point{x: 1, y: 1}
	p3 := &point{x: 0, y: 1}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p1, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{seg1, seg2, seg3})
	terr(t, err)

	ring := rings[0]
	ring.forceExteriorRing, ring.isExteriorRing = true, false

	expect(t, !ring.calcIsExteriorRing())
	expect(t, equalRing(ring.getGeom(), [][]float64{{0, 0}, {0, 1}, {1, 1}, {0, 0}}))
}

func TestGeomOutRingRemoveColinearPoints(t *testing.T) {
	t.Parallel()

	// removes colinear points successfully

	op := newOperation("")

	p1 := &point{x: 0, y: 0}
	p2 := &point{x: 1, y: 1}
	p3 := &point{x: 2, y: 2}
	p4 := &point{x: 0, y: 2}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p4, &ringIn{})
	terr(t, err)
	seg4, err := op.newSegmentFromRing(p4, p1, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	seg4.forceInResult, seg4.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{seg1, seg2, seg3, seg4})
	terr(t, err)

	ring := rings[0]

	expect(t, equalRing(ring.getGeom(), [][]float64{{0, 0}, {2, 2}, {0, 2}, {0, 0}}))
}

func TestGeomOutRingAlmostEqualPoint(t *testing.T) {
	t.Parallel()

	// almost equal point handled ok
	// points harvested from https://github.com/mfogel/polygon-clipping/issues/37

	op := newOperation("")

	p1 := &point{x: 0.523985, y: 51.281651}
	p2 := &point{x: 0.5241, y: 51.2816}
	p3 := &point{x: 0.5240213684210527, y: 51.2816873684210}
	p4 := &point{x: 0.5239850000000027, y: 51.281651000000004}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p4, &ringIn{})
	terr(t, err)
	seg4, err := op.newSegmentFromRing(p4, p1, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	seg4.forceInResult, seg4.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{seg1, seg2, seg3, seg4})
	terr(t, err)

	ring := rings[0]

	expect(t, equalRing(ring.getGeom(), [][]float64{
		{0.523985, 51.281651},
		{0.5241, 51.2816},
		{0.5240213684210527, 51.2816873684210},
		{0.523985, 51.281651}},
	))
}

func TestGeomOutRingWithColinearPoints(t *testing.T) {
	t.Parallel()

	// ring with all colinear points returns null

	op := newOperation("")

	p1 := &point{x: 0, y: 0}
	p2 := &point{x: 1, y: 1}
	p3 := &point{x: 2, y: 2}
	p4 := &point{x: 3, y: 3}

	seg1, err := op.newSegmentFromRing(p1, p2, &ringIn{})
	terr(t, err)
	seg2, err := op.newSegmentFromRing(p2, p3, &ringIn{})
	terr(t, err)
	seg3, err := op.newSegmentFromRing(p3, p4, &ringIn{})
	terr(t, err)
	seg4, err := op.newSegmentFromRing(p4, p1, &ringIn{})
	terr(t, err)

	seg1.forceInResult, seg1.inResult = true, true
	seg2.forceInResult, seg2.inResult = true, true
	seg3.forceInResult, seg3.inResult = true, true
	seg4.forceInResult, seg4.inResult = true, true

	rings, err := newRingOutFromSegments([]*segment{seg1, seg2, seg3, seg4})
	terr(t, err)

	ring := rings[0]

	expect(t, ring.getGeom() == nil)
}

// poly
func TestGeomOutPolyBasic(t *testing.T) {
	t.Parallel()

	// basic

	ring1 := &ringOut{
		forceGeom: true,
		geom:      [][]float64{{1}},
	}

	ring2 := &ringOut{
		forceGeom: true,
		geom:      [][]float64{{2}},
	}
	ring3 := &ringOut{
		forceGeom: true,
		geom:      [][]float64{{3}},
	}

	poly := newPolyOut(ring1)
	poly.addInterior(ring2)
	poly.addInterior(ring3)

	expect(t, ring1.poly == poly)
	expect(t, ring2.poly == poly)
	expect(t, ring3.poly == poly)

	expect(t, equalPoly(poly.getGeom(), [][][]float64{{{1}}, {{2}}, {{3}}}))
}

func TestGeomOutPolyColinearExteriorRing(t *testing.T) {
	t.Parallel()

	// has all colinear exterior ring

	ring1 := &ringOut{
		forceGeom: true,
		geom:      nil,
	}
	poly := newPolyOut(ring1)

	expect(t, ring1.poly == poly)
	expect(t, poly.getGeom() == nil)
}

func TestGeomOutPolyColinearInteriorRing(t *testing.T) {
	t.Parallel()

	// has all colinear interior ring

	ring1 := &ringOut{
		forceGeom: true,
		geom:      [][]float64{{1}},
	}
	ring2 := &ringOut{
		forceGeom: true,
		geom:      nil,
	}
	ring3 := &ringOut{
		forceGeom: true,
		geom:      [][]float64{{3}},
	}

	poly := newPolyOut(ring1)
	poly.addInterior(ring2)
	poly.addInterior(ring3)

	expect(t, ring1.poly == poly)
	expect(t, ring2.poly == poly)
	expect(t, ring3.poly == poly)

	expect(t, equalPoly(poly.getGeom(), [][][]float64{{{1}}, {{3}}}))
}

func TestGeomOutMultiPolyBasic(t *testing.T) {
	t.Parallel()

	// basic

	multiPoly := newMultiPolyOut([]*ringOut{})
	poly1 := &polyOut{
		forceGeom: true,
		geom:      [][][]float64{{{0}}},
	}
	poly2 := &polyOut{
		forceGeom: true,
		geom:      [][][]float64{{{1}}},
	}
	multiPoly.polys = []*polyOut{poly1, poly2}

	expect(t, equalMultiPoly(multiPoly.getGeom(), [][][][]float64{{{{0}}}, {{{1}}}}))
}

func TestGeomOutMultiPolyColinearExteriorRing(t *testing.T) {
	t.Parallel()

	// has poly with all colinear exterior ring

	multiPoly := newMultiPolyOut([]*ringOut{})
	poly1 := &polyOut{
		forceGeom: true,
		geom:      nil,
	}
	poly2 := &polyOut{
		forceGeom: true,
		geom:      [][][]float64{{{1}}},
	}
	multiPoly.polys = []*polyOut{poly1, poly2}

	expect(t, equalMultiPoly(multiPoly.getGeom(), [][][][]float64{{{{1}}}}))
}
