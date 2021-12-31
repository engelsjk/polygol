package polygol

import (
	"fmt"

	splaytree "github.com/engelsjk/splay-tree"
)

/**
 * NOTE:  We must be careful not to change any segments while
 *        they are in the SplayTree. AFAIK, there's no way to tell
 *        the tree to rebalance itself - thus before splitting
 *        a segment that's in the tree, we remove it from the tree,
 *        do the split, then re-insert it. (Even though splitting a
 *        segment *shouldn't* change its correct position in the
 *        sweep line tree, the reality is because of rounding errors,
 *        it sometimes does.)
 */

type sweepLine struct {
	tree     *splaytree.SplayTree
	queue    *splaytree.SplayTree
	segments []*segment
}

func newSweepLine(
	queue *splaytree.SplayTree,
	comparator func(a, b interface{}) int,
) *sweepLine {
	sl := &sweepLine{}
	if comparator == nil {
		comparator = segmentCompare
	}
	sl.queue = queue
	sl.tree = splaytree.New(comparator)
	sl.segments = []*segment{}
	return sl
}

func (sl *sweepLine) process(event *sweepEvent) ([]*sweepEvent, error) {

	seg := event.segment
	newEvents := []*sweepEvent{}

	// if we've already been consumed by another segment,
	// clean up our body parts and get out
	if event.consumedBy != nil {
		if event.isLeft {
			sl.queue.Remove(event.otherSE)
		} else {
			sl.tree.Remove(seg)
		}
		return newEvents, nil
	}

	var node *splaytree.Node
	if event.isLeft {
		node = sl.tree.Insert(seg)
	} else {
		node = sl.tree.Find(seg)
	}

	if node == nil {
		// fmt.Printf("event: %s\n", fmtSweepEvent(event))
		// fmt.Printf("segment: %s\n", fmtSegment(segment))
		// fmt.Printf("node: %+v\n", node)
		// logTree(sl.tree)

		return nil, fmt.Errorf(
			`Unable to find segment #%d [%f, %f] -> [%f, %f] in SweepLine tree. 
			Please submit a bug report.`,
			seg.id,
			seg.leftSE.point.x, seg.leftSE.point.y,
			seg.rightSE.point.x, seg.rightSE.point.y,
		)
	}

	prevNode := node
	nextNode := node
	var prevSeg *segment
	var nextSeg *segment

	// skip consumed segments still in tree
	for prevSeg == nil {
		prevNode = sl.tree.Prev(prevNode)
		if prevNode == nil {
			break
		} else if prevNode.Item().(*segment).consumedBy == nil {
			prevSeg = prevNode.Item().(*segment)
		}
	}

	// skip consumed segments still in tree
	for nextSeg == nil {
		nextNode = sl.tree.Next(nextNode)
		if nextNode == nil {
			break
		} else if nextNode.Item().(*segment).consumedBy == nil {
			nextSeg = nextNode.Item().(*segment)
		}
	}

	if event.isLeft {

		var prevSplitter, nextSplitter *point

		// Check for intersections against the previous segment in the sweep line
		prevSplitter, newEvents = sl.getSplitterFromIntersections(seg, prevSeg, newEvents)

		// Check for intersections against the next segment in the sweep line
		nextSplitter, newEvents = sl.getSplitterFromIntersections(seg, nextSeg, newEvents)

		// For simplicity, even if we find more than one intersection we only
		// spilt on the 'earliest' (sweep-line style) of the intersections.
		// The other intersection will be handled in a future process().
		if prevSplitter != nil || nextSplitter != nil {

			var splitter *point
			if prevSplitter == nil {
				splitter = nextSplitter
			} else if nextSplitter == nil {
				splitter = prevSplitter
			} else {
				cmpSplitters := sweepEventComparePoints(prevSplitter, nextSplitter)
				splitter = nextSplitter
				if cmpSplitters <= 0 {
					splitter = prevSplitter
				}
			}

			// Rounding errors can cause changes in ordering,
			// so remove affected segments and right sweep events before splitting
			sl.queue.Remove(seg.rightSE)
			newEvents = append(newEvents, seg.rightSE)

			newEventsFromSplit := seg.split(splitter)
			newEvents = append(newEvents, newEventsFromSplit...)
		}

		if len(newEvents) > 0 {
			// We found some intersections, so re-do the current event to
			// make sure sweep line ordering is totally consistent for later
			// use with the segment 'prev' pointers
			sl.tree.Remove(seg)
			newEvents = append(newEvents, event)

		} else {
			// done with left event
			sl.segments = append(sl.segments, seg)
			seg.prev = prevSeg
		}
	} else {
		// event.isRight

		// since we're about to be removed from the sweep line, check for
		// intersections between our previous and next segments

		if prevSeg != nil && nextSeg != nil {
			inter := prevSeg.getIntersection(nextSeg)
			if inter != nil {
				newEvents = sl.splitOnInter(prevSeg, inter, newEvents)
				newEvents = sl.splitOnInter(nextSeg, inter, newEvents)
			}
		}
		sl.tree.Remove(seg)
	}

	return newEvents, nil
}

// splitSafely splits a segment that is currently in the datastructures
// IE - a segment other than the one that is currently being processed.
func (sl *sweepLine) splitSafely(segment *segment, point *point) []*sweepEvent {
	// Rounding errors can cause changes in ordering,
	// so remove affected segments and right sweep events before splitting
	// removeNode() doesn't work, so have re-find the seg
	// https://github.com/w8r/splay-tree/pull/5
	sl.tree.Remove(segment)
	rightSE := segment.rightSE
	sl.queue.Remove(rightSE)
	newEvents := segment.split(point)
	newEvents = append(newEvents, rightSE)
	// splitting can trigger consumption
	if segment.consumedBy == nil {
		sl.tree.Insert(segment)
	}
	return newEvents
}

func (sl *sweepLine) getSplitterFromIntersections(seg, other *segment, events []*sweepEvent) (*point, []*sweepEvent) {

	var splitter *point
	if other != nil {
		otherInter := other.getIntersection(seg)
		if otherInter != nil {
			if !seg.isAndEndpoint(otherInter) {
				splitter = otherInter
			}
			if !other.isAndEndpoint(otherInter) {
				newEventsFromSplit := sl.splitSafely(other, otherInter)
				for i := 0; i < len(newEventsFromSplit); i++ {
					events = append(events, newEventsFromSplit[i])
				}
			}
		}
	}
	return splitter, events
}

func (sl *sweepLine) splitOnInter(seg *segment, inter *point, events []*sweepEvent) []*sweepEvent {
	if !seg.isAndEndpoint(inter) {
		newEventsFromSplit := sl.splitSafely(seg, inter)
		for i := 0; i < len(newEventsFromSplit); i++ {
			events = append(events, newEventsFromSplit[i])
		}
	}
	return events
}
