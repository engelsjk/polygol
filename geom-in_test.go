package polygol

import (
	"testing"
)

func TestRingIn(t *testing.T) {
	var ringIn *ringIn
	var err error
	var ring [][]float64

	op := newOperation("")

	// create exterior ring
	ring = [][]float64{
		{0, 0},
		{1, 0},
		{1, 1},
	}
	expectedPt1 := &point{x: 0, y: 0}
	expectedPt2 := &point{x: 1, y: 0}
	expectedPt3 := &point{x: 1, y: 1}
	poly := &polyIn{}
	ringIn, err = op.newRingIn(ring, poly, true)
	terr(t, err)
	poly.exteriorRing = ringIn

	expect(t, ringIn.poly == poly)
	expect(t, ringIn.isExterior)
	expect(t, len(ringIn.segments) == 3)
	expect(t, len(ringIn.getSweepEvents()) == 6)

	expect(t, ringIn.segments[0].leftSE.point.equal(*expectedPt1))
	expect(t, ringIn.segments[0].rightSE.point.equal(*expectedPt2))
	expect(t, ringIn.segments[1].leftSE.point.equal(*expectedPt2))
	expect(t, ringIn.segments[1].rightSE.point.equal(*expectedPt3))
	expect(t, ringIn.segments[2].leftSE.point.equal(*expectedPt1))
	expect(t, ringIn.segments[2].rightSE.point.equal(*expectedPt3))

	// create an interior ring
	ring = [][]float64{{0, 0}, {1, 1}, {1, 0}}
	ringIn, err = op.newRingIn(ring, &polyIn{}, false)
	terr(t, err)
	expect(t, !ringIn.isExterior)
}

func TestPolyIn(t *testing.T) {

	op := newOperation("")

	// creation
	multiPolyIn := &multiPolyIn{}
	poly := [][][]float64{
		{{0, 0}, {10, 0}, {10, 10}, {0, 10}},
		{{0, 0}, {1, 1}, {1, 0}},
		{{2, 2}, {2, 3}, {3, 3}, {3, 2}},
	}
	polyIn, err := op.newPolyIn(poly, multiPolyIn)
	terr(t, err)

	expect(t, polyIn.multiPoly == multiPolyIn)
	expect(t, len(polyIn.exteriorRing.segments) == 4)
	expect(t, len(polyIn.interiorRings) == 2)
	expect(t, len(polyIn.interiorRings[0].segments) == 3)
	expect(t, len(polyIn.interiorRings[1].segments) == 4)
	expect(t, len(polyIn.getSweepEvents()) == 22)
}

func TestmultiPolyIn(t *testing.T) {

	var multiPolyIn *multiPolyIn
	var err error

	op := newOperation("")

	// creation with multipoly
	multiPolyIn, err = op.newMultiPolyIn([][][][]float64{
		{{{0, 0}, {1, 1}, {0, 1}}},
		{
			{{0, 0}, {4, 0}, {4, 9}},
			{{2, 2}, {3, 3}, {3, 2}},
		},
	}, false)
	terr(t, err)

	expect(t, len(multiPolyIn.polys) == 2)
	expect(t, len(multiPolyIn.getSweepEvents()) == 18)

	// creation with poly
	multiPolyIn, err = op.newMultiPolyIn([][][][]float64{
		{{{0, 0}, {1, 1}, {0, 1}, {0, 0}}},
	}, false)
	terr(t, err)

	expect(t, len(multiPolyIn.polys) == 1)
	expect(t, len(multiPolyIn.getSweepEvents()) == 6)

	// third or more coordinates are ignored
	multiPolyIn, err = op.newMultiPolyIn([][][][]float64{
		{{{0, 0, 42}, {1, 1, 128}, {0, 1, 84}, {0, 0, 42}}},
	}, false)
	terr(t, err)

	expect(t, len(multiPolyIn.polys) == 1)
	expect(t, len(multiPolyIn.getSweepEvents()) == 6)

	// creation with invalid input
	// creation with point
	// creation with ring
	// creation with empty polygon / ring
	// creation with empty ring / point
	// creation with polygon with invalid coordinates
	// creation with polygon with missing coordinates
	// creation with multipolygon with invalid coordinates

	// creation with multipolygon with missing coordinates
	_, err = op.newMultiPolyIn([][][][]float64{
		{{{0}, {0, 1}, {1, 0}}},
	}, false)
	expect(t, err != nil)
}
