package polygol

type point struct {
	x      float64
	y      float64
	events []*sweepEvent
}

func newPoint(x, y float64) *point {
	return &point{
		x: x, y: y,
	}
}

func (p point) xy() []float64 {
	return []float64{p.x, p.y}
}

func (p point) equal(point point) bool {
	return p.x == point.x && p.y == point.y
}
