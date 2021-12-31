package polygol

import (
	"fmt"
	"sort"
)

type linkedEvent struct {
	index int
	point *point
}

type geomOut interface {
	geom() Geom
}

type ringOut struct {
	geom              [][]float64
	forceGeom         bool
	events            []*sweepEvent
	poly              *polyOut
	enclosingRing     *ringOut
	isExteriorRing    bool
	forceExteriorRing bool
}

func newRingOut(events []*sweepEvent) *ringOut {
	ro := &ringOut{}
	ro.events = events
	for i := 0; i < len(events); i++ {
		events[i].segment.ringOut = ro
	}
	return ro
}

func newRingOutFromSegments(allSegments []*segment) ([]*ringOut, error) {

	ringsOut := []*ringOut{}

	for i := 0; i < len(allSegments); i++ {
		segment := allSegments[i]
		if !segment.isInResult() || segment.ringOut != nil {
			continue
		}

		var prevEvent *sweepEvent
		event := segment.leftSE
		nextEvent := segment.rightSE
		events := []*sweepEvent{event}

		startingPoint := event.point
		intersectionLEs := []linkedEvent{}

		// Walk the chain of linked events to form a closed ring
		for {

			prevEvent = event
			event = nextEvent
			events = append(events, event)

			// Is the ring complete?
			if event.point == startingPoint {
				break
			}

			for {
				availableLEs := event.getAvailableLinkedEvents()

				// Did we hit a dead end? This shouldn't happen. Indicates some earlier
				//  part of the algorithm malfunctioned... please file a bug report.
				if len(availableLEs) == 0 {
					firstPt := events[0].point
					lastPt := events[len(events)-1].point
					return nil, fmt.Errorf(`Unable to complete output ring starting at [%f, %f].
					Last matching segment found ends at [%f, %f].`,
						firstPt.x, firstPt.y, lastPt.x, lastPt.y)
				}

				// Only one way to go, so continue on the path.
				if len(availableLEs) == 1 {
					nextEvent = availableLEs[0].otherSE
					break
				}

				// We must have an intersection. Check for a completed loop.
				indexLE := -1
				for j := 0; j < len(intersectionLEs); j++ {
					if intersectionLEs[j].point == event.point {
						indexLE = j
						break
					}
				}

				// Found a completed loop. Cut that off and make a ring.
				if indexLE >= 0 {
					intersectionLE := intersectionLEs[indexLE] // splice
					intersectionLEs = intersectionLEs[:indexLE]

					ringEvents := events[intersectionLE.index:] // splice w/ index
					events = events[:intersectionLE.index]

					ringEvents = append([]*sweepEvent{ringEvents[0].otherSE}, ringEvents...) // unshift
					for i, j := 0, len(ringEvents)-1; i < j; i, j = i+1, j-1 {               // reverse ringEvents
						ringEvents[i], ringEvents[j] = ringEvents[j], ringEvents[i]
					}
					ringsOut = append(ringsOut, newRingOut(ringEvents))
					continue
				}

				// Register the intersection.
				intersectionLEs = append(intersectionLEs, linkedEvent{index: len(events), point: event.point})

				// Choose the left-most option to continue the walk.
				comparator := event.getLeftMostComparator(prevEvent)
				sort.SliceStable(availableLEs, func(i, j int) bool {
					return comparator(availableLEs[i], availableLEs[j]) < 0
				})
				nextEvent = availableLEs[0].otherSE
				break
			}
		}
		ringsOut = append(ringsOut, newRingOut(events))
	}
	return ringsOut, nil
}

func (ro *ringOut) getGeom() [][]float64 {
	if ro.forceGeom {
		return ro.geom
	}
	// Remove superfluous points (ie extra points along a straight line),
	prevPt := ro.events[0].point
	points := []*point{prevPt}
	for i := 1; i < len(ro.events)-1; i++ {
		pt := ro.events[i].point
		nextPt := ro.events[i+1].point
		if compareAngles(
			[]float64{pt.x, pt.y},
			[]float64{prevPt.x, prevPt.y},
			[]float64{nextPt.x, nextPt.y},
		) == 0 {
			continue
		}
		points = append(points, pt)
		prevPt = pt
	}

	// ring was all (within rounding error of angle calc) colinear points
	if len(points) == 1 {
		return nil
	}

	// check if the starting point is necessary
	pt := points[0]
	nextPt := points[1]
	if compareAngles(
		[]float64{pt.x, pt.y},
		[]float64{prevPt.x, prevPt.y},
		[]float64{nextPt.x, nextPt.y},
	) == 0 {
		points = points[1:]
	}

	points = append(points, points[0])
	step := -1
	if ro.calcIsExteriorRing() {
		step = 1
	}
	iStart := len(points) - 1
	if ro.calcIsExteriorRing() {
		iStart = 0
	}
	iEnd := -1
	if ro.calcIsExteriorRing() {
		iEnd = len(points)
	}
	orderedPoints := [][]float64{}
	for i := iStart; i != iEnd; i += step {
		orderedPoints = append(orderedPoints, []float64{points[i].x, points[i].y})
	}
	return orderedPoints
}

func (ro *ringOut) calcIsExteriorRing() bool {
	if ro.forceExteriorRing {
		return ro.isExteriorRing
	}
	if ro.enclosingRing == nil {
		enclosing := ro.getEnclosingRing()
		if enclosing != nil {
			ro.isExteriorRing = !enclosing.calcIsExteriorRing()
		} else {
			ro.isExteriorRing = true
		}
	}
	return ro.isExteriorRing
}

func (ro *ringOut) getEnclosingRing() *ringOut {
	if ro.enclosingRing == nil {
		ro.enclosingRing = ro.calcEnclosingRing()
	}
	return ro.enclosingRing
}

func (ro *ringOut) calcEnclosingRing() *ringOut {
	// start with the earlier sweep line event so that the prevSeg
	// chain doesn't lead us inside of a loop of ours
	leftMostEvt := ro.events[0]
	for i := 1; i < len(ro.events); i++ {
		evt := ro.events[i]
		if sweepEventCompare(leftMostEvt, evt) > 0 {
			leftMostEvt = evt
		}
	}

	prevSeg := leftMostEvt.segment.prevInResult()
	var prevPrevSeg *segment
	if prevSeg != nil {
		prevPrevSeg = prevSeg.prevInResult()
	}

	for {
		// no segment found, thus no ring can enclose us
		if prevSeg == nil {
			return nil
		}

		// no segments below prev segment found, thus the ring of the prev
		// segment must loop back around and enclose us
		if prevPrevSeg == nil {
			return prevSeg.ringOut
		}

		// if the two segments are of different rings, the ring of the prev
		// segment must either loop around us or the ring of the prev prev
		// seg, which would make us and the ring of the prev peers
		if prevPrevSeg.ringOut != prevSeg.ringOut {
			if prevPrevSeg.ringOut.getEnclosingRing() != prevSeg.ringOut {
				return prevSeg.ringOut
			} else {
				return prevSeg.ringOut.getEnclosingRing()
			}
		}

		// two segments are from the same ring, so this was a penisula
		// of that ring. iterate downward, keep searching
		prevSeg = prevPrevSeg.prevInResult()
		prevPrevSeg = nil
		if prevSeg != nil {
			prevPrevSeg = prevSeg.prevInResult()
		}
	}
}

type polyOut struct {
	init          bool
	geom          [][][]float64
	forceGeom     bool
	exteriorRing  *ringOut
	poly          *polyOut
	interiorRings []*ringOut
}

func newPolyOut(exteriorRing *ringOut) *polyOut {
	po := &polyOut{}
	po.exteriorRing = exteriorRing
	exteriorRing.poly = po
	po.interiorRings = []*ringOut{}
	return po
}

func (po *polyOut) addInterior(ring *ringOut) {
	po.interiorRings = append(po.interiorRings, ring)
	ring.poly = po
}

func (po *polyOut) getGeom() [][][]float64 {
	if po.forceGeom {
		return po.geom
	}
	geom := [][][]float64{po.exteriorRing.getGeom()}
	// exterior ring was all (within rounding error of angle calc) colinear points
	if geom == nil {
		return nil
	}
	if geom[0] == nil {
		return nil
	}
	for i := 0; i < len(po.interiorRings); i++ {
		ringGeom := po.interiorRings[i].getGeom()
		// interior ring was all (within rounding error of angle calc) colinear points
		if ringGeom == nil {
			continue
		}
		geom = append(geom, ringGeom)
	}
	return geom
}

type multiPolyOut struct {
	geom      [][][][]float64
	forceGeom bool
	rings     []*ringOut
	polys     []*polyOut
}

func NewMultiPolyOut(rings []*ringOut) multiPolyOut {
	mpo := multiPolyOut{}
	mpo.rings = rings
	mpo.polys = mpo.composePolys(rings)
	return mpo
}

func (mpo *multiPolyOut) getGeom() [][][][]float64 {
	if mpo.forceGeom {
		return mpo.geom
	}
	geom := [][][][]float64{}
	for i := 0; i < len(mpo.polys); i++ {
		polyGeom := mpo.polys[i].getGeom()
		// exterior ring was all (within rounding error of angle calc) colinear points
		if polyGeom == nil {
			continue
		}
		geom = append(geom, polyGeom)
	}
	return geom
}

func (mpo *multiPolyOut) composePolys(rings []*ringOut) []*polyOut {
	polys := []*polyOut{}
	for i := 0; i < len(rings); i++ {
		ring := rings[i]
		if ring.poly != nil {
			continue
		}
		if ring.calcIsExteriorRing() {
			polys = append(polys, newPolyOut(ring))
		} else {
			enclosingRing := ring.getEnclosingRing()
			if enclosingRing.poly == nil {
				polys = append(polys, newPolyOut(enclosingRing))
			}
			enclosingRing.poly.addInterior(ring)
		}
	}
	return polys
}

//////////////
