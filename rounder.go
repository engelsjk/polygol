package polygol

import (
	splaytree "github.com/engelsjk/splay-tree"
)

type ptRounder struct {
	xRounder *coordRounder
	yRounder *coordRounder
}

func newPtRounder() *ptRounder {
	ptr := new(ptRounder)
	ptr.reset()
	return ptr
}

func (pr *ptRounder) reset() {
	pr.xRounder = newCoordRounder()
	pr.yRounder = newCoordRounder()
}

func (pr *ptRounder) round(x, y float64) *point {
	return newPoint(
		pr.xRounder.round(x),
		pr.yRounder.round(y),
	)
}

type coordRounder struct {
	tree *splaytree.SplayTree
}

func newCoordRounder() *coordRounder {
	cr := new(coordRounder)
	less := func(a, b interface{}) int {
		af := a.(float64)
		bf := b.(float64)
		if af > bf {
			return 1
		}
		if af < bf {
			return -1
		}
		return 0
	}
	cr.tree = splaytree.New(less)
	cr.round(0.0)
	return cr
}

func (cr *coordRounder) round(coord float64) float64 {

	node := cr.tree.Add(coord)
	item := node.Item().(float64)

	prevNode := cr.tree.Prev(node)
	if prevNode != nil {
		prevItem := prevNode.Item().(float64)
		if flpCmp(item, prevItem) == 0 {
			cr.tree.Remove(coord)
			return prevItem
		}
	}

	nextNode := cr.tree.Next(node)
	if nextNode != nil {
		nextItem := nextNode.Item().(float64)
		if flpCmp(item, nextItem) == 0 {
			cr.tree.Remove(coord)
			return nextItem
		}
	}

	return coord
}
