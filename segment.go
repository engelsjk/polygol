package polygol

import (
	"fmt"
)

var segmentID = 1

type segment struct {
	init            bool
	id              int
	leftSE          *sweepEvent
	rightSE         *sweepEvent
	rings           []*ringIn
	windings        []int
	ringOut         *ringOut
	consumedBy      *segment
	inResult        bool
	forceInResult   bool
	doneInResult    bool
	prev            *segment
	prevSegInResult *segment
	after           *state
	before          *state
}

func newSegment(leftSE, rightSE *sweepEvent, rings []*ringIn, windings []int) *segment {
	s := &segment{}
	segmentID++
	s.id = segmentID
	s.leftSE = leftSE
	leftSE.segment = s
	leftSE.otherSE = rightSE
	s.rightSE = rightSE
	rightSE.segment = s
	rightSE.otherSE = leftSE
	s.rings = rings
	s.windings = windings
	s.init = true
	return s
}

func segmentCompare(a, b interface{}) int {

	aSeg := a.(*segment)
	bSeg := b.(*segment)

	alx := aSeg.leftSE.point.x
	blx := bSeg.leftSE.point.x
	arx := aSeg.rightSE.point.x
	brx := bSeg.rightSE.point.x

	// check if they're even in the same vertical plane
	if brx < alx {
		return 1
	}
	if arx < blx {
		return -1
	}

	aly := aSeg.leftSE.point.y
	bly := bSeg.leftSE.point.y
	ary := aSeg.rightSE.point.y
	bry := bSeg.rightSE.point.y

	// is left endpoint of segment B the right-more?
	if alx < blx {

		// are the two segments in the same horizontal plane?
		if bly < aly && bly < ary {
			return 1
		}
		if bly > aly && bly > ary {
			return -1
		}

		// is the B left endpoint colinear to segment A?
		aCmpBLeft := aSeg.comparePoint(bSeg.leftSE.point)
		if aCmpBLeft < 0 {
			return 1
		}
		if aCmpBLeft > 0 {
			return -1
		}

		// is the A right endpoint colinear to segment B ?
		bCmpARight := bSeg.comparePoint(aSeg.rightSE.point)
		if bCmpARight != 0 {
			return bCmpARight
		}

		// colinear segments, consider the one with left-more
		// left endpoint to be first (arbitrary?)
		return -1
	}

	// is left endpoint of segment A the right-more?
	if alx > blx {

		if aly < bly && aly < bry {
			return -1
		}
		if aly > bly && aly > bry {
			return 1
		}

		// is the A left endpoint colinear to segment B?
		bCmpALeft := bSeg.comparePoint(aSeg.leftSE.point)
		if bCmpALeft != 0 {
			return bCmpALeft
		}

		// is the B right endpoint colinear to segment A?
		aCmpBRight := aSeg.comparePoint(bSeg.rightSE.point)
		if aCmpBRight < 0 {
			return 1
		}
		if aCmpBRight > 0 {
			return -1
		}

		// colinear segments, consider the one with left-more
		// left endpoint to be first (arbitrary?)
		return 1
	}

	// if we get here, the two left endpoints are in the same
	// vertical plane, ie alx === blx

	// consider the lower left-endpoint to come first
	if aly < bly {
		return -1
	}
	if aly > bly {
		return 1
	}

	// left endpoints are identical
	// check for colinearity by using the left-more right endpoint

	// is the A right endpoint more left-more?
	if arx < brx {
		bCmpARight := bSeg.comparePoint(aSeg.rightSE.point)
		if bCmpARight != 0 {
			return bCmpARight
		}
	}

	// is the B right endpoint more left-more?
	if arx > brx {
		aCmpBRight := aSeg.comparePoint(bSeg.rightSE.point)
		if aCmpBRight < 0 {
			return 1
		}
		if aCmpBRight > 0 {
			return -1
		}
	}

	if arx != brx {
		// are these two [almost] vertical segments with opposite orientation?
		// if so, the one with the lower right endpoint comes first
		ay := ary - aly
		ax := arx - alx
		by := bry - bly
		bx := brx - blx
		if ay > ax && by < bx {
			return 1
		}
		if ay < ax && by > bx {
			return -1
		}
	}

	// we have colinear segments with matching orientation
	// consider the one with more left-more right endpoint to be first
	if arx > brx {
		return 1
	}
	if arx < brx {
		return -1
	}

	// if we get here, two two right endpoints are in the same
	// vertical plane, ie arx === brx

	// consider the lower right-endpoint to come first
	if ary < bry {
		return -1
	}
	if ary > bry {
		return 1
	}

	// right endpoints identical as well, so the segments are identical
	// fall back on creation order as consistent tie-breaker
	if aSeg.id < bSeg.id {
		return -1
	}
	if aSeg.id > bSeg.id {
		return 1
	}

	// identical segment, ie a === b
	return 0
}

func newSegmentFromRing(pt1, pt2 *point, ring *ringIn) (*segment, error) {

	var leftPt, rightPt *point
	var winding int

	cmpPts := sweepEventComparePoints(pt1, pt2)
	if cmpPts < 0 {
		leftPt = pt1
		rightPt = pt2
		winding = 1
	} else if cmpPts > 0 {
		leftPt = pt2
		rightPt = pt1
		winding = -1
	} else {
		return nil, fmt.Errorf("Tried to create degenerate segment at [%f,%f].", pt1.x, pt1.y)
	}

	leftSE := newSweepEvent(leftPt, true)
	rightSE := newSweepEvent(rightPt, false)

	return newSegment(leftSE, rightSE, []*ringIn{ring}, []int{winding}), nil
}

func (s *segment) replaceRightSE(newRightSE *sweepEvent) {
	s.rightSE = newRightSE
	s.rightSE.segment = s
	s.rightSE.otherSE = s.leftSE
	s.leftSE.otherSE = s.rightSE
}

func (s *segment) bbox() bbox {

	y1 := s.leftSE.point.y
	y2 := s.rightSE.point.y

	lly := y2
	if y1 < y2 {
		lly = y1
	}

	ury := y2
	if y1 > y2 {
		ury = y1
	}

	return bbox{
		ll: point{x: s.leftSE.point.x, y: lly},
		ur: point{x: s.rightSE.point.x, y: ury},
	}
}

func (s *segment) vector() []float64 {
	return []float64{
		s.rightSE.point.x - s.leftSE.point.x,
		s.rightSE.point.y - s.leftSE.point.y,
	}
}

func (s *segment) isAndEndpoint(point *point) bool {
	if s == nil {
		return false
	}
	if point == nil {
		return false
	}
	if s.leftSE == nil {
		return false
	}
	return (point.x == s.leftSE.point.x && point.y == s.leftSE.point.y) ||
		(point.x == s.rightSE.point.x && point.y == s.rightSE.point.y)
}

func (s *segment) comparePoint(point *point) int {

	if s.isAndEndpoint(point) {
		return 0
	}

	lPt := s.leftSE.point
	rPt := s.rightSE.point
	v := s.vector()

	// Exactly vertical segments.
	if lPt.x == rPt.x {
		if point.x == lPt.x {
			return 0
		}
		if point.x < lPt.x {
			return 1
		}
		return -1
	}

	// Nearly vertical segments with an intersection.
	// Check to see where a point on the line with matching Y coordinate is.

	yDist := (point.y - lPt.y) / v[1]
	xFromYDist := lPt.x + yDist*v[0]
	if point.x == xFromYDist {
		return 0
	}

	// General case.
	// Check to see where a point on the line with matching X coordinate is.

	xDist := (point.x - lPt.x) / v[0]
	yFromXDist := lPt.y + xDist*v[1]
	if point.y == yFromXDist {
		return 0
	}
	if point.y < yFromXDist {
		return -1
	}
	return 1
}

func (s *segment) getIntersection(other *segment) *point {

	if s == nil {
		return nil
	}
	if other == nil {
		return nil
	}

	// If bboxes don't overlap, there can't be any intersections
	tBbox := s.bbox()
	oBbox := other.bbox()

	bboxOverlap := tBbox.getBboxOverlap(oBbox)
	if bboxOverlap == nil {
		return nil
	}

	// We first check to see if the endpoints can be considered intersections.
	// This will 'snap' intersections to endpoints if possible, and will
	// handle cases of colinearity.

	tlp := s.leftSE.point
	trp := s.rightSE.point
	olp := other.leftSE.point
	orp := other.rightSE.point

	// does each endpoint touch the other segment?
	// note that we restrict the 'touching' definition to only allow segments
	// to touch endpoints that lie forward from where we are in the sweep line pass
	touchesOtherLSE := tBbox.isInBbox(*olp) && s.comparePoint(olp) == 0
	touchesThisLSE := oBbox.isInBbox(*tlp) && other.comparePoint(tlp) == 0
	touchesOtherRSE := tBbox.isInBbox(*orp) && s.comparePoint(orp) == 0
	touchesThisRSE := oBbox.isInBbox(*trp) && other.comparePoint(trp) == 0

	// do left endpoints match?
	if touchesThisLSE && touchesOtherLSE {
		// these two cases are for colinear segments with matching left
		// endpoints, and one segment being longer than the other
		if touchesThisRSE && !touchesOtherRSE {
			return trp
		}
		if !touchesThisRSE && touchesOtherRSE {
			return orp
		}
		// either the two segments match exactly (two trival intersections)
		// or just on their left endpoint (one trivial intersection
		return nil
	}

	// does this left endpoint matches (other doesn't)
	if touchesThisLSE {
		// check for segments that just intersect on opposing endpoints
		if touchesOtherRSE {
			if tlp.x == orp.x && tlp.y == orp.y {
				return nil
			}
		}
		// t-intersection on left endpoint
		return tlp
	}

	// does other left endpoint matches (this doesn't)
	if touchesOtherLSE {
		// check for segments that just intersect on opposing endpoints
		if touchesThisRSE {
			if trp.x == olp.x && trp.y == olp.y {
				return nil
			}
		}
		// t-intersection on left endpoint
		return olp
	}

	// trivial intersection on right endpoints
	if touchesThisRSE && touchesOtherRSE {
		return nil
	}

	// t-intersections on just one right endpoint
	if touchesThisRSE {
		return trp
	}
	if touchesOtherRSE {
		return orp
	}

	// None of our endpoints intersect. Look for a general intersection between
	// infinite lines laid over the segments

	pt := intersection(
		s.vector(),
		other.vector(),
		tlp.xy(),
		olp.xy(),
	)
	var ptInter *point
	if pt != nil {
		ptInter = newPoint(pt[0], pt[1])
	}

	// ptInter := lineToLineIntersection(
	// 	s.leftSE.point, s.rightSE.point,
	// 	other.leftSE.point, other.rightSE.point)

	defer func() { ptInter = nil }() // clean up dangling pointer

	// are the segments parallel? Note that if they were colinear with overlap,
	// they would have an endpoint intersection and that case was already handled above
	if ptInter == nil {
		return nil
	}

	// is the intersection found between the lines not on the segments?
	if !bboxOverlap.isInBbox(*ptInter) {
		return nil
	}

	return rounder.round(ptInter.x, ptInter.y)
}

func lineToLineIntersection(
	line1Start, line1End,
	line2Start, line2End *point,
) *point {
	// from github.com/twpayne/go-geom/xy/lineintersector/nonrobust_line_intersector.go

	var a2, b2 float64
	var c2, r1, r2, r3, r4 float64

	a1 := line1End.y - line1Start.y
	b1 := line1Start.x - line1End.x
	c1 := line1End.x*line1Start.y - line1Start.x*line1End.y

	r3 = a1*line2Start.x + b1*line2Start.y + c1
	r4 = a1*line2End.x + b1*line2End.y + c1

	if r3 != 0 && r4 != 0 && isSameSignAndNonZero(r3, r4) {
		return nil
	}

	a2 = line2End.y - line2Start.y
	b2 = line2Start.x - line2End.x
	c2 = line2End.x*line2Start.y - line2Start.x*line2End.y

	r1 = a2*line1Start.x + b2*line1Start.y + c2
	r2 = a2*line1End.x + b2*line1End.y + c2

	if r1 != 0 && r2 != 0 && isSameSignAndNonZero(r1, r2) {
		return nil
	}

	denom := a1*b2 - a2*b1
	if denom == 0 {
		/// ??? collinear intersection?
		return nil
	}

	numX := b1*c2 - b2*c1
	numY := a2*c1 - a1*c2

	return newPoint(numX/denom, numY/denom)
}

func (s *segment) split(point *point) []*sweepEvent {

	newEvents := []*sweepEvent{}
	alreadyLinked := point.events != nil

	newLeftSE := newSweepEvent(point, true)
	newRightSE := newSweepEvent(point, false)
	oldRightSE := s.rightSE

	s.replaceRightSE(newRightSE)
	newEvents = append(newEvents, newRightSE)
	newEvents = append(newEvents, newLeftSE)

	newRings := make([]*ringIn, len(s.rings))
	copy(newRings, s.rings)

	newWindings := make([]int, len(s.windings))
	copy(newWindings, s.windings)

	newSeg := newSegment(newLeftSE, oldRightSE, newRings, newWindings)

	// newSeg := NewSegment(newLeftSE, oldRightSE, s.rings, s.windings)

	// when splitting a nearly vertical downward-facing segment,
	// sometimes one of the resulting new segments is vertical, in which
	// case its left and right events may need to be swapped
	if sweepEventComparePoints(newSeg.leftSE.point, newSeg.rightSE.point) > 0 {
		newSeg.swapEvents()
	}
	if sweepEventComparePoints(s.leftSE.point, s.rightSE.point) > 0 {
		s.swapEvents()
	}

	// in the point we just used to create new sweep events with was already
	// linked to other events, we need to check if either of the affected
	// segments should be consumed
	if alreadyLinked {
		newLeftSE.checkForConsuming()
		newRightSE.checkForConsuming()
	}

	return newEvents
}

func (s *segment) swapEvents() {
	tmpEvt := s.rightSE
	s.rightSE = s.leftSE
	s.leftSE = tmpEvt
	s.leftSE.isLeft = true
	s.rightSE.isLeft = false
	for i := 0; i < len(s.windings); i++ {
		s.windings[i] *= -1
	}
}

func (s *segment) consume(other *segment) {

	consumer := s
	consumee := other

	for consumer.consumedBy != nil {
		consumer = consumer.consumedBy
	}
	for consumee.consumedBy != nil {
		consumee = consumee.consumedBy
	}

	cmp := segmentCompare(consumer, consumee)
	if cmp == 0 {
		return // already consumed
	}

	// the winner of the consumption is the earlier segment
	// according to sweep line ordering
	if cmp > 0 {
		tmp := consumer
		consumer = consumee
		consumee = tmp
	}

	// make sure a segment doesn't consume its prev
	if consumer.prev == consumee {
		tmp := consumer
		consumer = consumee
		consumee = tmp
	}

	for i := 0; i < len(consumee.rings); i++ {
		ring := consumee.rings[i]
		winding := consumee.windings[i]
		index := ring.indexOf(consumer.rings)
		if index == -1 {
			consumer.rings = append(consumer.rings, ring)
			consumer.windings = append(consumer.windings, winding)
		} else {
			consumer.windings[index] += winding
		}
	}
	consumee.rings = nil
	consumee.windings = nil
	consumee.consumedBy = consumer

	// mark sweep events consumed as to maintain ordering in sweep event queue
	consumee.leftSE.consumedBy = consumer.leftSE
	consumee.rightSE.consumedBy = consumer.rightSE
}

func (s *segment) prevInResult() *segment {
	if s.prevSegInResult != nil {
		return s.prevSegInResult
	}
	if s.prev == nil {
		s.prevSegInResult = nil
	} else if s.prev.isInResult() {
		s.prevSegInResult = s.prev
	} else {
		s.prevSegInResult = s.prev.prevInResult()
	}
	return s.prevSegInResult
}

type state struct {
	rings      []*ringIn
	windings   []int
	multiPolys []*multiPolyIn
}

func (s *segment) beforeState() *state {

	if s.before != nil {
		return s.before
	}
	if s.prev == nil {
		s.before = &state{
			rings:      []*ringIn{},
			windings:   []int{},
			multiPolys: []*multiPolyIn{},
		}
	} else {
		seg := s.prev.consumedBy
		if s.prev.consumedBy == nil {
			seg = s.prev
		}
		s.before = seg.afterState()
	}
	return s.before
}

func (s *segment) afterState() *state {

	if s.after != nil {
		return s.after
	}

	beforeState := s.beforeState()

	ringsBefore := make([]*ringIn, len(beforeState.rings))
	copy(ringsBefore, beforeState.rings)
	windingsBefore := make([]int, len(beforeState.windings))
	copy(windingsBefore, beforeState.windings)

	s.after = &state{
		rings:      ringsBefore,
		windings:   windingsBefore,
		multiPolys: []*multiPolyIn{},
	}

	// calculate ringsAfter, windingsAfter
	for i := 0; i < len(s.rings); i++ {
		ring := s.rings[i]
		winding := s.windings[i]
		index := ring.indexOf(s.after.rings)
		if index == -1 {
			s.after.rings = append(s.after.rings, ring)
			s.after.windings = append(s.after.windings, winding)
		} else {
			s.after.windings[index] += winding
		}
	}

	// calculate polysAfter
	polysAfter := []*polyIn{}
	polysExclude := []*polyIn{}
	for i := 0; i < len(s.after.rings); i++ {
		if s.after.windings[i] == 0 { // non-zero rule
			continue
		}
		ring := s.after.rings[i]
		poly := ring.poly
		index := poly.indexOf(polysExclude)
		if index != -1 {
			continue
		}
		if ring.isExterior { // exterior ring
			polysAfter = append(polysAfter, poly)
		} else { // interior ring
			if poly.indexOf(polysExclude) == -1 {
				polysExclude = append(polysExclude, poly)
			}
			index := ring.poly.indexOf(polysAfter)
			if index != -1 {
				polysAfter = append(polysAfter[:index], polysAfter[index+1:]...) // splice
			}
		}
	}

	// calculate multiPolysAfter
	for i := 0; i < len(polysAfter); i++ {
		mp := polysAfter[i].multiPoly
		if mp.indexOf(s.after.multiPolys) == -1 {
			s.after.multiPolys = append(s.after.multiPolys, mp)
		}
	}

	return s.after
}

func (s *segment) isInResult() bool {
	// if we've been consumed, we're not in the result
	if s == nil {
		return false
	}
	if s.consumedBy != nil {
		return false
	}
	if s.forceInResult {
		return s.inResult
	}
	if s.doneInResult {
		return s.inResult
	}

	mpsBefore := s.beforeState().multiPolys
	mpsAfter := s.afterState().multiPolys

	switch operationType {
	case "union":
		// UNION - included iff:
		//  * On one side of us there is 0 poly interiors AND
		//  * On the other side there is 1 or more.
		noBefores := len(mpsBefore) == 0
		noAfters := len(mpsAfter) == 0
		s.inResult = noBefores != noAfters
		s.doneInResult = true
	case "intersection":
		// INTERSECTION - included iff:
		//  * on one side of us all multipolys are rep. with poly interiors AND
		//  * on the other side of us, not all multipolys are repsented
		//    with poly interiors
		var least, most int
		if len(mpsBefore) < len(mpsAfter) {
			least = len(mpsBefore)
			most = len(mpsAfter)
		} else {
			least = len(mpsAfter)
			most = len(mpsBefore)
		}
		s.inResult = most == op.numMultiPolys && least < most
		s.doneInResult = true
	case "xor":
		// XOR - included iff:
		//  * the difference between the number of multipolys represented
		//    with poly interiors on our two sides is an odd number
		diff := abs(len(mpsBefore) - len(mpsAfter))
		s.inResult = diff%2 == 1
		s.doneInResult = true
	case "difference":
		// DIFFERENCE included iff:
		//  * on exactly one side, we have just the subject
		isJustSubject := func(mps []*multiPolyIn) bool {
			return len(mps) == 1 && mps[0].isSubject
		}
		s.inResult = isJustSubject(mpsBefore) != isJustSubject(mpsAfter)
		s.doneInResult = true
	default:
		fmt.Printf("Unrecognized operation type found %s", operationType)
	}
	return s.inResult
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func isSameSignAndNonZero(a, b float64) bool {
	if a == 0 || b == 0 {
		return false
	}
	return (a < 0 && b < 0) || (a > 0 && b > 0)
}
