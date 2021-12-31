package polygol

type bbox struct {
	ll point
	ur point
}

func (b bbox) isInBbox(point point) bool {
	return b.ll.x <= point.x &&
		point.x <= b.ur.x &&
		b.ll.y <= point.y &&
		point.y <= b.ur.y
}

func (b bbox) getBboxOverlap(ob bbox) *bbox {
	// check if the bboxes overlap at all
	if ob.ur.x < b.ll.x ||
		b.ur.x < ob.ll.x ||
		ob.ur.y < b.ll.y ||
		b.ur.y < ob.ll.y {
		return nil
	}

	// find the middle two X values
	lowerX := b.ll.x
	if b.ll.x < ob.ll.x {
		lowerX = ob.ll.x
	}
	upperX := ob.ur.x
	if b.ur.x < ob.ur.x {
		upperX = b.ur.x
	}

	// find the middle two Y values
	lowerY := b.ll.y
	if b.ll.y < ob.ll.y {
		lowerY = ob.ll.y
	}
	upperY := ob.ur.y
	if b.ur.y < ob.ur.y {
		upperY = b.ur.y
	}

	return &bbox{
		ll: point{x: lowerX, y: lowerY},
		ur: point{x: upperX, y: upperY},
	}
}
