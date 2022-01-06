package polygol

import (
	"testing"
)

func TestGeomInRingIn(t *testing.T) {
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

func TestGeomInPolyIn(t *testing.T) {

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

func TestGeomInMultiPolyIn(t *testing.T) {

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

func TestGeomInRingInIndexOf(t *testing.T) {
	r1 := &ringIn{}
	r2 := &ringIn{}
	r3 := &ringIn{}
	r4 := &ringIn{}
	r5 := &ringIn{}
	r6 := &ringIn{}

	ringIns := []*ringIn{r1, r2, r3, r4, r5}
	expect(t, r1.indexOf(ringIns) == 0)
	expect(t, r2.indexOf(ringIns) == 1)
	expect(t, r3.indexOf(ringIns) == 2)
	expect(t, r4.indexOf(ringIns) == 3)
	expect(t, r5.indexOf(ringIns) == 4)
	expect(t, r6.indexOf(ringIns) == -1)
}

func TestGeomInPolyInIndexOf(t *testing.T) {
	p1 := &polyIn{}
	p2 := &polyIn{}
	p3 := &polyIn{}
	p4 := &polyIn{}
	p5 := &polyIn{}
	p6 := &polyIn{}

	polyIns := []*polyIn{p1, p2, p3, p4, p5}
	expect(t, p1.indexOf(polyIns) == 0)
	expect(t, p2.indexOf(polyIns) == 1)
	expect(t, p3.indexOf(polyIns) == 2)
	expect(t, p4.indexOf(polyIns) == 3)
	expect(t, p5.indexOf(polyIns) == 4)
	expect(t, p6.indexOf(polyIns) == -1)
}

func TestGeomInMultiPolyInIndexOf(t *testing.T) {
	mp1 := &multiPolyIn{}
	mp2 := &multiPolyIn{}
	mp3 := &multiPolyIn{}
	mp4 := &multiPolyIn{}
	mp5 := &multiPolyIn{}
	mp6 := &multiPolyIn{}

	multiPolyIns := []*multiPolyIn{mp1, mp2, mp3, mp4, mp5}
	expect(t, mp1.indexOf(multiPolyIns) == 0)
	expect(t, mp2.indexOf(multiPolyIns) == 1)
	expect(t, mp3.indexOf(multiPolyIns) == 2)
	expect(t, mp4.indexOf(multiPolyIns) == 3)
	expect(t, mp5.indexOf(multiPolyIns) == 4)
	expect(t, mp6.indexOf(multiPolyIns) == -1)
}
