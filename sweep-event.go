package polygol

import (
	"errors"
	"fmt"
)

type sweepEvent struct {
	point      *point
	isLeft     bool
	segment    *segment
	consumedBy *sweepEvent
	otherSE    *sweepEvent
}

type angles struct {
	sine, cosine float64
}

func newSweepEvent(point *point, isLeft bool) *sweepEvent {
	se := &sweepEvent{}
	if point.events == nil {
		(*point).events = []*sweepEvent{se}
	} else {
		(*point).events = append((*point).events, se)
	}
	se.point = point
	se.isLeft = isLeft

	return se
}

// Compare orders sweep events in the sweep event queue
func sweepEventCompare(a, b interface{}) int {

	aSE := a.(*sweepEvent)
	bSE := b.(*sweepEvent)

	ptCmp := sweepEventComparePoints(aSE.point, bSE.point)
	if ptCmp != 0 {
		return ptCmp
	}

	// the points are the same, so link them if needed
	if aSE.point != bSE.point {
		if err := aSE.link(bSE); err != nil {
			fmt.Println(err)
		}
	}

	// favor right events over left
	if aSE.isLeft != bSE.isLeft {
		if aSE.isLeft {
			return 1
		}
		return -1
	}

	// we have two matching left or right endpoints
	// ordering of this case is the same as for their segments
	return segmentCompare(aSE.segment, bSE.segment)
}

func sweepEventComparePoints(aPt, bPt *point) int {
	cmpX := flpCmp(aPt.x, bPt.x)
	if cmpX != 0 {
		return cmpX
	}
	return flpCmp(aPt.y, bPt.y)
}

func (se *sweepEvent) link(other *sweepEvent) error {
	if other.point == se.point {
		return errors.New("Tried to link already linked events.")
	}
	otherEvents := other.point.events
	for i := 0; i < len(otherEvents); i++ {
		evt := otherEvents[i]
		se.point.events = append(se.point.events, evt)
		evt.point = se.point
	}
	se.checkForConsuming()
	return nil
}

// Do a pass over our linked events and check to see if any pair
// of segments match, and should be consumed.
func (se *sweepEvent) checkForConsuming() {
	// TODO: The loops in this method run O(n^2) => no good.
	//        Maintain little ordered sweep event trees?
	//        Can we maintaining an ordering that avoids the need
	//        for the re-sorting with getLeftmostComparator in geom-out?

	// Compare each pair of events to see if other events also match

	numEvents := len(se.point.events)

	for i := 0; i < numEvents; i++ {
		evt1 := se.point.events[i]
		if evt1.segment.consumedBy != nil {
			continue
		}
		for j := i + 1; j < numEvents; j++ {
			evt2 := se.point.events[j]
			if evt2.segment.consumedBy != nil {
				continue
			}

			if !equalSweepEvents(evt1.otherSE.point.events, evt2.otherSE.point.events) { // more correct? or not
				// if &(evt1.otherSE.point.events) != &(evt2.otherSE.point.events) {
				continue
			}
			evt1.segment.consume(evt2.segment)
		}
	}
}

func (se *sweepEvent) getAvailableLinkedEvents() []*sweepEvent {
	// point.events is always of length 2 or greater
	events := []*sweepEvent{}
	for i := 0; i < len(se.point.events); i++ {
		evt := se.point.events[i]
		isInResult := evt.segment.isInResult()
		if !isInResult {
			continue
		}
		if se != evt && evt.segment.ringOut == nil {
			events = append(events, evt)
		}
	}
	return events
}

func (se *sweepEvent) getLeftMostComparator(baseEvent *sweepEvent) func(a, b *sweepEvent) int {
	cache := make(map[*sweepEvent]angles)
	fillCache := func(linkedEvent *sweepEvent) {
		nextEvent := linkedEvent.otherSE
		cache[linkedEvent] = angles{
			sine:   sineOfAngle(se.point.xy(), baseEvent.point.xy(), nextEvent.point.xy()),
			cosine: cosineOfAngle(se.point.xy(), baseEvent.point.xy(), nextEvent.point.xy()),
		}
	}

	return func(a, b *sweepEvent) int {
		if _, ok := cache[a]; !ok {
			fillCache(a)
		}
		if _, ok := cache[b]; !ok {
			fillCache(b)
		}
		aa := cache[a]
		bb := cache[b]

		// both on or above x-axis
		if aa.sine >= 0 && bb.sine >= 0 {
			if aa.cosine < bb.cosine {
				return 1
			}
			if aa.cosine > bb.cosine {
				return -1
			}
			return 0
		}

		// both below x-axis
		if aa.sine < 0 && bb.sine < 0 {
			if aa.cosine < bb.cosine {
				return -1
			}
			if aa.cosine > bb.cosine {
				return 1
			}
			return 0
		}

		// one above x-axis, one below
		if bb.sine < aa.sine {
			return -1
		}
		if bb.sine > aa.sine {
			return 1
		}
		return 0
	}
}

func equalSweepEvents(a, b []*sweepEvent) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
