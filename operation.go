package polygol

import (
	"fmt"

	splaytree "github.com/engelsjk/splay-tree"
)

var (
	polygonClippingMaxQueueSize         = 1000000
	polygonClippingMaxSweepLineSegments = 1000000
)

type Geom [][][][]float64

var operationType string

type operation struct {
	numMultiPolys int
}

func (op *operation) run(optype string, geom Geom, moreGeoms ...Geom) (Geom, error) {

	operationType = optype
	rounder.reset()

	multiPolys, err := geomsToMultiPolys(geom, moreGeoms)
	if err != nil {
		return Geom{}, err
	}
	op.numMultiPolys = len(multiPolys)

	switch operationType {
	// BBox optimization for difference operation
	// If the bbox of a multipolygon that's part of the clipping doesn't
	// intersect the bbox of the subject at all, we can just drop that
	// multipolygon.
	case "difference":
		// in place removal
		subject := multiPolys[0]
		i := 1
		for i < len(multiPolys) {
			if multiPolys[i].bbox.getBboxOverlap(subject.bbox) != nil {
				i++
			} else {
				multiPolys = append(multiPolys[:i], multiPolys[i+1:]...) // splice
			}
		}
	// BBox optimization for intersection operation
	// If we can find any pair of multipolygons whose bbox does not overlap,
	// then the result will be empty.
	case "intersection":
		// TODO: this is O(n^2) in number of polygons. By sorting the bboxes,
		//       it could be optimized to O(n * ln(n))
		for i := 0; i < len(multiPolys); i++ {
			mpA := multiPolys[i]
			for j := i + 1; j < len(multiPolys); j++ {
				if mpA.bbox.getBboxOverlap(multiPolys[i].bbox) == nil {
					return Geom{}, nil
				}
			}
		}
	}

	// Put segment endpoints in a priority queue.
	// Should be sorted by x coordinate.
	queue := splaytree.New(sweepEventCompare)
	for i := 0; i < len(multiPolys); i++ {
		sweepEvents := multiPolys[i].getSweepEvents()
		for j := 0; j < len(sweepEvents); j++ {
			queue.Insert(sweepEvents[j])
			if queue.Size() > polygonClippingMaxQueueSize {
				// prevents an infinite loop, an otherwise common manifestation of bugs
				return nil, fmt.Errorf(`Infinite loop when putting segment endpoints in a priority queue 
				(queue size too big). Please file a bug report.`)
			}
		}
	}

	///////////////////////////////////////////////////////////////
	///////////////////////////////////////////////////////////////

	// Pass the sweep line over those endpoints.
	sweepLine := newSweepLine(queue, nil)
	prevQueueSize := queue.Size()
	node := queue.Pop()
	i := 0
	for node != nil {

		evt := node.Item().(*sweepEvent)
		if queue.Size() == prevQueueSize {

			// prevents an infinite loop, an otherwise common manifestation of bugs
			seg := evt.segment
			dir := "right"
			if evt.isLeft {
				dir = "left"
			}
			return nil, fmt.Errorf(`Unable to pop() %s SweepEvent [%f, %f]
			from segment #%d [%f, %f] -> [%f, %f] from queue. Please file a bug report.`,
				dir,
				evt.point.x, evt.point.y,
				seg.id,
				seg.leftSE.point.x, seg.leftSE.point.y,
				seg.rightSE.point.x, seg.rightSE.point.y)
		}

		if queue.Size() > polygonClippingMaxQueueSize {
			// prevents an infinite loop, an otherwise common manifestation of bugs
			return nil, fmt.Errorf(`Infinite loop when passing sweep line over endspoints
			(queue size too big). Please file a bug report.`)
		}

		if len(sweepLine.segments) > polygonClippingMaxSweepLineSegments {
			// prevents an infinite loop, an otherwise common manifestation of bugs
			return nil, fmt.Errorf(`Infinite loop when passing sweep line over endspoints
			(too many sweep line segments). Please file a bug report.`)
		}

		newEvents, err := sweepLine.process(evt)
		if err != nil {
			return nil, err
		}

		for i := 0; i < len(newEvents); i++ {
			evt := newEvents[i]
			if evt.consumedBy == nil {
				queue.Insert(evt)
			}
		}
		prevQueueSize = queue.Size()
		node = queue.Pop()
		i++
	}

	// Free some memory we don't need anymore.
	rounder.reset()

	// Collect and compile segments we're keeping into a multipolygon.
	ringsOut, err := newRingOutFromSegments(sweepLine.segments)
	if err != nil {
		return nil, err
	}

	result := NewMultiPolyOut(ringsOut)

	return result.getGeom(), nil
}

func geomsToMultiPolys(geom Geom, moreGeoms []Geom) ([]*multiPolyIn, error) {
	// Convert inputs to MultiPoly objects.
	multiPoly, err := newMultiPolyIn(geom, true)
	if err != nil {
		return nil, err
	}
	multiPolys := []*multiPolyIn{multiPoly}
	for i := 0; i < len(moreGeoms); i++ {
		multiPoly, err := newMultiPolyIn(moreGeoms[i], false)
		if err != nil {
			continue
		}
		multiPolys = append(multiPolys, multiPoly)
	}
	return multiPolys, nil
}
