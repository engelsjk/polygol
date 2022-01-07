package polygol

import (
	"fmt"
	"os"
	"strconv"

	splaytree "github.com/engelsjk/splay-tree"
)

var (
	polygolClippingMaxQueueSize         = 1000000
	polygolClippingMaxSweepLineSegments = 1000000
)

func init() {
	envMaxQueueSize := os.Getenv("POLYGOL_MAX_QUEUE_SIZE")
	if envMaxQueueSize != "" {
		maxQueueSize, err := strconv.Atoi(envMaxQueueSize)
		if err != nil {
			fmt.Println("env var POLYGOL_MAX_QUEUE_SIZE must be a integer")
		}
		polygolClippingMaxQueueSize = maxQueueSize
	}
	envMaxSweepLineSegments := os.Getenv("POLYGOL_MAX_SWEEPLINE_SEGMENTS")
	if envMaxSweepLineSegments != "" {
		maxSweepLineSegments, err := strconv.Atoi(envMaxSweepLineSegments)
		if err != nil {
			fmt.Println("env var POLYGOL_MAX_SWEEPLINE_SEGMENTS must be a integer")
		}
		polygolClippingMaxSweepLineSegments = maxSweepLineSegments
	}
}

type operation struct {
	rounder       *ptRounder
	opType        string
	numMultiPolys int
	segmentID     int
}

func newOperation(opType string) *operation {
	rounder := newPtRounder()
	return &operation{
		rounder: rounder,
		opType:  opType,
	}
}

func (o *operation) run(geom Geom, moreGeoms ...Geom) (Geom, error) {

	o.rounder.reset()

	multiPolys, err := o.geomsToMultiPolys(geom, moreGeoms)
	if err != nil {
		return Geom{}, err
	}
	o.numMultiPolys = len(multiPolys)

	switch o.opType {
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
			if queue.Size() > polygolClippingMaxQueueSize {
				// prevents an infinite loop, an otherwise common manifestation of bugs
				return nil, fmt.Errorf(`Infinite loop when putting segment endpoints in a priority queue 
				(queue size too big). Try increasing POLYGOL_MAX_QUEUE_SIZE.`)
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

		if queue.Size() > polygolClippingMaxQueueSize {
			// prevents an infinite loop, an otherwise common manifestation of bugs
			return nil, fmt.Errorf(`Infinite loop when passing sweep line over endspoints
			(queue size too big). Try increasing POLYGOL_MAX_QUEUE_SIZE.`)
		}

		if len(sweepLine.segments) > polygolClippingMaxSweepLineSegments {
			// prevents an infinite loop, an otherwise common manifestation of bugs
			return nil, fmt.Errorf(`Infinite loop when passing sweep line over endspoints
			(too many sweep line segments). Try increasing POLYGOL_MAX_SWEEPLINE_SEGMENTS.`)
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
	o.rounder.reset()

	// Collect and compile segments we're keeping into a multipolygon.
	ringsOut, err := newRingOutFromSegments(sweepLine.segments)
	if err != nil {
		return nil, err
	}

	result := newMultiPolyOut(ringsOut)

	return result.getGeom(), nil
}

func (o *operation) geomsToMultiPolys(geom Geom, moreGeoms []Geom) ([]*multiPolyIn, error) {
	// Convert inputs to MultiPoly objects.
	multiPoly, err := o.newMultiPolyIn(geom, true)
	if err != nil {
		return nil, err
	}
	multiPolys := []*multiPolyIn{multiPoly}
	for i := 0; i < len(moreGeoms); i++ {
		multiPoly, err := o.newMultiPolyIn(moreGeoms[i], false)
		if err != nil {
			continue
		}
		multiPolys = append(multiPolys, multiPoly)
	}
	return multiPolys, nil
}
