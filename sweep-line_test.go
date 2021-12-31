package polygol

import "testing"

var comparator = func(a, b interface{}) int {
	af := a.(float64)
	bf := b.(float64)
	if af == bf {
		return 0
	}
	if af < bf {
		return -1
	}
	return 1
}

func TestSweepLine(t *testing.T) {

	// test filling up the tree then emptying it out
	sl := newSweepLine(nil, comparator)
	k1 := 4.0
	k2 := 9.0
	k3 := 13.0
	k4 := 44.0

	n1 := sl.tree.Insert(k1)
	n2 := sl.tree.Insert(k2)
	n4 := sl.tree.Insert(k4)
	n3 := sl.tree.Insert(k3)

	expect(t, sl.tree.Find(k1) == n1)
	expect(t, sl.tree.Find(k2) == n2)
	expect(t, sl.tree.Find(k3) == n3)
	expect(t, sl.tree.Find(k4) == n4)

	expect(t, sl.tree.Prev(n1) == nil)
	expect(t, sl.tree.Next(n1).Item() == k2)

	expect(t, sl.tree.Prev(n2).Item() == k1)
	expect(t, sl.tree.Next(n2).Item() == k3)

	expect(t, sl.tree.Prev(n3).Item() == k2)
	expect(t, sl.tree.Next(n3).Item() == k4)

	expect(t, sl.tree.Prev(n4).Item() == k3)
	expect(t, sl.tree.Next(n4) == nil)

	sl.tree.Remove(k2)
	expect(t, sl.tree.Find(k2) == nil)

	n1 = sl.tree.Find(k1)
	n3 = sl.tree.Find(k3)
	n4 = sl.tree.Find(k4)

	expect(t, sl.tree.Prev(n1) == nil)
	expect(t, sl.tree.Next(n1).Item() == k3)

	expect(t, sl.tree.Prev(n3).Item() == k1)
	expect(t, sl.tree.Next(n3).Item() == k4)

	expect(t, sl.tree.Prev(n4).Item() == k3)
	expect(t, sl.tree.Next(n4) == nil)

	sl.tree.Remove(k4)
	expect(t, sl.tree.Find(k4) == nil)

	n1 = sl.tree.Find(k1)
	n3 = sl.tree.Find(k3)

	expect(t, sl.tree.Prev(n1) == nil)
	expect(t, sl.tree.Next(n1).Item() == k3)

	expect(t, sl.tree.Prev(n3).Item() == k1)
	expect(t, sl.tree.Next(n3) == nil)

	sl.tree.Remove(k1)
	expect(t, sl.tree.Find(k1) == nil)

	n3 = sl.tree.Find(k3)

	expect(t, sl.tree.Prev(n3) == nil)
	expect(t, sl.tree.Next(n3) == nil)

	sl.tree.Remove(k3)
	expect(t, sl.tree.Find(k3) == nil)
}
