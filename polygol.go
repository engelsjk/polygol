package polygol

var rounder *ptRounder
var op operation

func init() {
	rounder = newPtRounder()
}

func Union(geom Geom, moreGeoms ...Geom) (Geom, error) {
	return op.run("union", geom, moreGeoms...)
}

func Intersection(geom Geom, moreGeoms ...Geom) (Geom, error) {
	return op.run("intersection", geom, moreGeoms...)
}

func XOR(geom Geom, moreGeoms ...Geom) (Geom, error) {
	return op.run("xor", geom, moreGeoms...)
}

func Difference(geom Geom, moreGeoms ...Geom) (Geom, error) {
	return op.run("difference", geom, moreGeoms...)
}
