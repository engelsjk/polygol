package polygol

import (
	"fmt"
	"math"
)

type ringIn struct {
	poly       *polyIn
	isExterior bool
	segments   []*segment
	bbox       bbox
}

func newRingIn(ring [][]float64, poly *polyIn, isExterior bool) (*ringIn, error) {

	if len(ring) == 0 {
		return nil, fmt.Errorf(`Input geometry is not a valid polygon or multipolygon (empty).`)
	}
	if len(ring[0]) < 2 {
		return nil, fmt.Errorf(`Input geometry is not a valid polygon or multipolygon (empty).`)
	}

	ri := &ringIn{}

	ri.poly = poly
	ri.isExterior = isExterior
	ri.segments = []*segment{}

	firstPoint := rounder.round(ring[0][0], ring[0][1])

	ri.bbox = bbox{ll: *firstPoint, ur: *firstPoint}

	prevPoint := firstPoint
	for i := 1; i < len(ring); i++ {

		if len(ring[i]) < 2 {
			return nil, fmt.Errorf(`Input geometry is not a valid polygon or multipolygon (missing coordinates).`)
		}

		point := rounder.round(ring[i][0], ring[i][1])

		// skip repeated points
		if point.x == prevPoint.x && point.y == prevPoint.y {
			continue
		}

		segment, err := newSegmentFromRing(prevPoint, point, ri)
		if err != nil {
			return nil, err
		}
		ri.segments = append(ri.segments, segment)

		if point.x < ri.bbox.ll.x {
			ri.bbox.ll.x = point.x
		}
		if point.y < ri.bbox.ll.y {
			ri.bbox.ll.y = point.y
		}
		if point.x > ri.bbox.ur.x {
			ri.bbox.ur.x = point.x
		}
		if point.y > ri.bbox.ur.y {
			ri.bbox.ur.y = point.y
		}
		prevPoint = point
	}
	// add segment from last to first if last is not the same as first
	if firstPoint.x != prevPoint.x || firstPoint.y != prevPoint.y {
		segment, err := newSegmentFromRing(prevPoint, firstPoint, ri)
		if err != nil {
			return nil, err
		}
		ri.segments = append(ri.segments, segment)
	}
	return ri, nil
}

func (ri *ringIn) getSweepEvents() []*sweepEvent {
	sweepEvents := []*sweepEvent{}
	for i := 0; i < len(ri.segments); i++ {
		segment := ri.segments[i]
		sweepEvents = append(sweepEvents, segment.leftSE, segment.rightSE)
	}
	return sweepEvents
}

func (ri *ringIn) indexOf(ringIns []*ringIn) int {
	for i, r := range ringIns {
		if ri == nil || r == nil {
			continue
		}
		if ri == r {
			return i
		}
	}
	return -1
}

type polyIn struct {
	multiPoly     *multiPolyIn
	exteriorRing  *ringIn
	interiorRings []*ringIn
	bbox          bbox
}

func newPolyIn(poly [][][]float64, multiPoly *multiPolyIn) (*polyIn, error) {

	if len(poly) == 0 {
		return nil, fmt.Errorf(`Input geometry is not a valid polygon or multipolygon (empty).`)
	}

	pi := &polyIn{}

	exteriorRing, err := newRingIn(poly[0], pi, true)
	if err != nil {
		return nil, err
	}

	pi.exteriorRing = exteriorRing
	pi.bbox = pi.exteriorRing.bbox
	pi.interiorRings = []*ringIn{}

	for i := 1; i < len(poly); i++ {
		ring, err := newRingIn(poly[i], pi, false)
		if err != nil {
			return nil, err
		}
		if ring.bbox.ll.x < pi.bbox.ll.x {
			pi.bbox.ll.x = ring.bbox.ll.x
		}
		if ring.bbox.ll.y < pi.bbox.ll.y {
			pi.bbox.ll.y = ring.bbox.ll.y
		}
		if ring.bbox.ur.x > pi.bbox.ur.x {
			pi.bbox.ur.x = ring.bbox.ur.x
		}
		if ring.bbox.ur.y > pi.bbox.ur.y {
			pi.bbox.ur.y = ring.bbox.ur.y
		}
		pi.interiorRings = append(pi.interiorRings, ring)
	}
	pi.multiPoly = multiPoly
	return pi, nil
}

func (pi *polyIn) getSweepEvents() []*sweepEvent {
	sweepEvents := pi.exteriorRing.getSweepEvents()
	for i := 0; i < len(pi.interiorRings); i++ {
		ringSweepEvents := pi.interiorRings[i].getSweepEvents()
		sweepEvents = append(sweepEvents, ringSweepEvents...)
	}
	return sweepEvents
}

func (pi *polyIn) indexOf(polyIns []*polyIn) int {
	for i, p := range polyIns {
		if pi == nil || p == nil {
			continue
		}
		if pi == p {
			return i
		}
	}
	return -1
}

type multiPolyIn struct {
	polys     []*polyIn
	bbox      bbox
	isSubject bool
}

func newMultiPolyIn(multiPoly [][][][]float64, isSubject bool) (*multiPolyIn, error) {

	mpi := &multiPolyIn{}

	mpi.polys = []*polyIn{}
	mpi.bbox = bbox{
		ll: point{x: math.Inf(1), y: math.Inf(1)},
		ur: point{x: math.Inf(-1), y: math.Inf(-1)},
	}

	for i := 0; i < len(multiPoly); i++ {
		poly, err := newPolyIn(multiPoly[i], mpi)
		if err != nil {
			return nil, err
		}
		if poly.bbox.ll.x < mpi.bbox.ll.x {
			mpi.bbox.ll.x = poly.bbox.ll.x
		}
		if poly.bbox.ll.y < mpi.bbox.ll.y {
			mpi.bbox.ll.y = poly.bbox.ll.y
		}
		if poly.bbox.ur.x > mpi.bbox.ur.x {
			mpi.bbox.ur.x = poly.bbox.ur.x
		}
		if poly.bbox.ur.y > mpi.bbox.ur.y {
			mpi.bbox.ur.y = poly.bbox.ur.y
		}
		mpi.polys = append(mpi.polys, poly)
	}
	mpi.isSubject = isSubject
	return mpi, nil
}

func (mpi *multiPolyIn) getSweepEvents() []*sweepEvent {
	sweepEvents := []*sweepEvent{}
	for i := 0; i < len(mpi.polys); i++ {
		polySweepEvents := mpi.polys[i].getSweepEvents()
		sweepEvents = append(sweepEvents, polySweepEvents...)
	}
	return sweepEvents
}

func (mpi *multiPolyIn) indexOf(multiPolyIns []*multiPolyIn) int {
	for i, mp := range multiPolyIns {
		if mpi == nil || mp == nil {
			continue
		}
		if mpi == mp {
			return i
		}
	}
	return -1
}
