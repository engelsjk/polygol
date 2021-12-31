package polygol

import (
	"testing"
)

func equalBbox(bb1, bb2 bbox) bool {
	return bb1.ll.x == bb2.ll.x &&
		bb1.ll.y == bb2.ll.y &&
		bb1.ur.x == bb2.ur.x &&
		bb1.ur.y == bb2.ur.y
}

func TestisInBbox(t *testing.T) {
	var b bbox

	// outside
	b = bbox{ll: point{x: 1, y: 2}, ur: point{x: 5, y: 6}}
	expect(t, !b.isInBbox(point{x: 0, y: 3}))
	expect(t, !b.isInBbox(point{x: 3, y: 30}))
	expect(t, !b.isInBbox(point{x: 3, y: -30}))
	expect(t, !b.isInBbox(point{x: 9, y: 3}))

	// inside
	b = bbox{ll: point{x: 1, y: 2}, ur: point{x: 5, y: 6}}
	expect(t, b.isInBbox(point{x: 1, y: 2}))
	expect(t, b.isInBbox(point{x: 5, y: 6}))
	expect(t, b.isInBbox(point{x: 1, y: 6}))
	expect(t, b.isInBbox(point{x: 5, y: 2}))
	expect(t, b.isInBbox(point{x: 3, y: 4}))

	// barely inside & outside
	b = bbox{ll: point{x: 1, y: 0.8}, ur: point{x: 1.2, y: 6}}
	expect(t, b.isInBbox(point{x: 1.2 - epsilon, y: 6}))
	expect(t, !b.isInBbox(point{x: 1.2 + epsilon, y: 6}))
	expect(t, b.isInBbox(point{x: 1, y: 0.8 + epsilon}))
	expect(t, !b.isInBbox(point{x: 1, y: 0.8 - epsilon}))
}

func TestBboxOverlap(t *testing.T) {

	var b1, b2 bbox
	b1 = bbox{ll: point{x: 4, y: 4}, ur: point{x: 6, y: 6}}

	// disjoint - none
	// above
	b2 = bbox{ll: point{x: 7, y: 7}, ur: point{x: 8, y: 8}}
	expect(t, b1.getBboxOverlap(b2) == nil)
	// left
	b2 = bbox{ll: point{x: 1, y: 5}, ur: point{x: 3, y: 8}}
	expect(t, b1.getBboxOverlap(b2) == nil)
	// down
	b2 = bbox{ll: point{x: 2, y: 2}, ur: point{x: 3, y: 3}}
	expect(t, b1.getBboxOverlap(b2) == nil)
	// right
	b2 = bbox{ll: point{x: 12, y: 1}, ur: point{x: 14, y: 9}}
	expect(t, b1.getBboxOverlap(b2) == nil)

	// touching - one point
	// upper right corner of 1
	b2 = bbox{ll: point{x: 6, y: 6}, ur: point{x: 7, y: 8}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 6, y: 6}, ur: point{x: 6, y: 6}}))
	// upper left corner of 1
	b2 = bbox{ll: point{x: 3, y: 6}, ur: point{x: 4, y: 8}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 4, y: 6}, ur: point{x: 4, y: 6}}))
	// lower left corner of 1
	b2 = bbox{ll: point{x: 0, y: 0}, ur: point{x: 4, y: 4}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 4, y: 4}, ur: point{x: 4, y: 4}}))
	// lower right corner of 1
	b2 = bbox{ll: point{x: 6, y: 0}, ur: point{x: 12, y: 4}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 6, y: 4}, ur: point{x: 6, y: 4}}))

	// overlapping - two points

	// full overlap

	// matching bboxes
	expect(t, equalBbox(*b1.getBboxOverlap(b1), b1))

	// one side & two corners matching
	b2 = bbox{ll: point{x: 4, y: 4}, ur: point{x: 5, y: 6}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 4, y: 4}, ur: point{x: 5, y: 6}}))

	// one corner matching, part of two sides
	b2 = bbox{ll: point{x: 5, y: 4}, ur: point{x: 6, y: 5}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 5, y: 4}, ur: point{x: 6, y: 5}}))

	// part of a side matching, no corners
	b2 = bbox{ll: point{x: 4.5, y: 4.5}, ur: point{x: 5.5, y: 6}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 4.5, y: 4.5}, ur: point{x: 5.5, y: 6}}))

	// completely enclosed - no side or corner matching
	b2 = bbox{ll: point{x: 4.5, y: 5}, ur: point{x: 5.5, y: 5.5}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), b2))

	// partial overlap

	// full side overlap
	b2 = bbox{ll: point{x: 3, y: 4}, ur: point{x: 5, y: 6}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 4, y: 4}, ur: point{x: 5, y: 6}}))

	// partial side overlap
	b2 = bbox{ll: point{x: 5, y: 4.5}, ur: point{x: 7, y: 5.5}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 5, y: 4.5}, ur: point{x: 6, y: 5.5}}))

	// corner overlap
	b2 = bbox{ll: point{x: 5, y: 5}, ur: point{x: 7, y: 7}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 5, y: 5}, ur: point{x: 6, y: 6}}))

	// line bboxes

	// vertical line & normal

	// no overlap
	b2 = bbox{ll: point{x: 7, y: 3}, ur: point{x: 7, y: 6}}
	expect(t, b1.getBboxOverlap(b2) == nil)

	// point overlap
	b2 = bbox{ll: point{x: 6, y: 0}, ur: point{x: 6, y: 4}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 6, y: 4}, ur: point{x: 6, y: 4}}))

	// line overlap
	b2 = bbox{ll: point{x: 5, y: 0}, ur: point{x: 5, y: 9}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 5, y: 4}, ur: point{x: 5, y: 6}}))

	// horizontal line & normal

	// no overlap
	b2 = bbox{ll: point{x: 3, y: 7}, ur: point{x: 6, y: 7}}
	expect(t, b1.getBboxOverlap(b2) == nil)

	// point overlap
	b2 = bbox{ll: point{x: 1, y: 6}, ur: point{x: 4, y: 6}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 4, y: 6}, ur: point{x: 4, y: 6}}))

	// line overlap
	b2 = bbox{ll: point{x: 4, y: 6}, ur: point{x: 6, y: 6}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 4, y: 6}, ur: point{x: 6, y: 6}}))

	// two vertical lines
	var v1, v2 bbox
	v1 = bbox{ll: point{x: 4, y: 4}, ur: point{x: 4, y: 6}}

	// no overlap
	v2 = bbox{ll: point{x: 4, y: 7}, ur: point{x: 4, y: 8}}
	expect(t, v1.getBboxOverlap(v2) == nil)

	// point overlap
	v2 = bbox{ll: point{x: 4, y: 3}, ur: point{x: 4, y: 4}}
	expect(t, equalBbox(*v1.getBboxOverlap(v2), bbox{ll: point{x: 4, y: 4}, ur: point{x: 4, y: 4}}))

	// line overlap
	v2 = bbox{ll: point{x: 4, y: 3}, ur: point{x: 4, y: 5}}
	expect(t, equalBbox(*v1.getBboxOverlap(v2), bbox{ll: point{x: 4, y: 4}, ur: point{x: 4, y: 5}}))

	// two horizontal lines
	var h1, h2 bbox
	h1 = bbox{ll: point{x: 4, y: 6}, ur: point{x: 7, y: 6}}

	// no overlap
	h2 = bbox{ll: point{x: 4, y: 5}, ur: point{x: 7, y: 5}}
	expect(t, h1.getBboxOverlap(h2) == nil)

	// point overlap
	h2 = bbox{ll: point{x: 7, y: 6}, ur: point{x: 8, y: 6}}
	expect(t, equalBbox(*h1.getBboxOverlap(h2), bbox{ll: point{x: 7, y: 6}, ur: point{x: 7, y: 6}}))

	// line overlap
	h2 = bbox{ll: point{x: 4, y: 6}, ur: point{x: 7, y: 6}}
	expect(t, equalBbox(*h1.getBboxOverlap(h2), bbox{ll: point{x: 4, y: 6}, ur: point{x: 7, y: 6}}))

	// horizontal & vertical lines

	// no overlap
	h1 = bbox{ll: point{x: 4, y: 6}, ur: point{x: 8, y: 6}}
	v1 = bbox{ll: point{x: 5, y: 7}, ur: point{x: 5, y: 9}}
	expect(t, h1.getBboxOverlap(v1) == nil)

	// point overlap
	h1 = bbox{ll: point{x: 4, y: 6}, ur: point{x: 8, y: 6}}
	v1 = bbox{ll: point{x: 5, y: 5}, ur: point{x: 5, y: 9}}
	expect(t, equalBbox(*h1.getBboxOverlap(v1), bbox{ll: point{x: 5, y: 6}, ur: point{x: 5, y: 6}}))

	// produced line box

	// horizontal
	b2 = bbox{ll: point{x: 4, y: 6}, ur: point{x: 8, y: 8}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 4, y: 6}, ur: point{x: 6, y: 6}}))

	// vertical
	b2 = bbox{ll: point{x: 6, y: 2}, ur: point{x: 8, y: 8}}
	expect(t, equalBbox(*b1.getBboxOverlap(b2), bbox{ll: point{x: 6, y: 4}, ur: point{x: 6, y: 6}}))

	// point bboxes
	var p bbox

	// point & normal

	// no overlap
	p = bbox{ll: point{x: 2, y: 2}, ur: point{x: 2, y: 2}}
	expect(t, b1.getBboxOverlap(p) == nil)

	// point overlap
	p = bbox{ll: point{x: 5, y: 5}, ur: point{x: 5, y: 5}}
	expect(t, equalBbox(*b1.getBboxOverlap(p), p))

	// point & line
	var l bbox

	// no overlap
	p = bbox{ll: point{x: 2, y: 2}, ur: point{x: 2, y: 2}}
	l = bbox{ll: point{x: 4, y: 6}, ur: point{x: 4, y: 8}}
	expect(t, l.getBboxOverlap(p) == nil)

	// point overlap
	p = bbox{ll: point{x: 5, y: 5}, ur: point{x: 5, y: 5}}
	l = bbox{ll: point{x: 4, y: 5}, ur: point{x: 6, y: 5}}
	expect(t, equalBbox(*l.getBboxOverlap(p), p))

	// point & point
	var p1, p2 bbox

	// no overlap
	p1 = bbox{ll: point{x: 2, y: 2}, ur: point{x: 2, y: 2}}
	p2 = bbox{ll: point{x: 4, y: 6}, ur: point{x: 4, y: 6}}
	expect(t, p1.getBboxOverlap(p2) == nil)

	// point overlap
	p = bbox{ll: point{x: 5, y: 5}, ur: point{x: 5, y: 5}}
	expect(t, equalBbox(*p.getBboxOverlap(p), p))
}
