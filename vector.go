package polygol

import (
	"math"
)

// func Intersection(v1, v2 []float64, pt1, pt2 []float64) []float64 {
func intersection(v1, v2 []float64, pt1, pt2 []float64) []float64 {
	// take some shortcuts for vertical and horizontal lines
	// this also ensures we don't calculate an intersection and then discover
	// it's actually outside the bounding box of the line
	if v1[0] == 0 {
		return verticalIntersection(v2, pt2, pt1[0])
	}
	if v2[0] == 0 {
		return verticalIntersection(v1, pt1, pt2[0])
	}
	if v1[1] == 0 {
		return horizontalIntersection(v2, pt2, pt1[1])
	}
	if v2[1] == 0 {
		return horizontalIntersection(v1, pt1, pt2[1])
	}

	// General case for non-overlapping segments.
	// This algorithm is based on Schneider and Eberly.
	// http://www.cimec.org.ar/~ncalvo/Schneider_Eberly.pdf - pg 244

	kross := crossProduct(v1, v2)
	if kross == 0 {
		return nil
	}

	ve := []float64{pt2[0] - pt1[0], pt2[1] - pt1[1]}
	d1 := crossProduct(ve, v1) / kross
	d2 := crossProduct(ve, v2) / kross

	// take the average of the two calculations to minimize rounding error
	x1, x2 := pt1[0]+d2*v1[0], pt2[0]+d1*v2[0]
	y1, y2 := pt1[1]+d2*v1[1], pt2[1]+d1*v2[1]
	x := (x1 + x2) / 2.0
	y := (y1 + y2) / 2.0
	return []float64{x, y}
}

func compareAngles(basePt, endPt1, endPt2 []float64) int {
	v1 := []float64{endPt1[0] - basePt[0], endPt1[1] - basePt[1]}
	v2 := []float64{endPt2[0] - basePt[0], endPt2[1] - basePt[1]}
	kross := crossProduct(v1, v2)
	return flpCmp(kross, 0)
}

func sineOfAngle(pShared, pBase, pAngle []float64) float64 {
	vBase := []float64{pBase[0] - pShared[0], pBase[1] - pShared[1]}
	vAngle := []float64{pAngle[0] - pShared[0], pAngle[1] - pShared[1]}
	return crossProduct(vAngle, vBase) / length(vAngle) / length(vBase)
}

func cosineOfAngle(pShared, pBase, pAngle []float64) float64 {
	vBase := []float64{pBase[0] - pShared[0], pBase[1] - pShared[1]}
	vAngle := []float64{pAngle[0] - pShared[0], pAngle[1] - pShared[1]}
	return dotProduct(vAngle, vBase) / length(vAngle) / length(vBase)
}

func length(v []float64) float64 {
	return math.Sqrt(dotProduct(v, v))
}

func equal(v1, v2 []float64) bool {
	return v1[0] == v2[0] && v1[1] == v2[1]
}

func crossProduct(v1, v2 []float64) float64 {
	return v1[0]*v2[1] - v1[1]*v2[0]
}

func dotProduct(v1, v2 []float64) float64 {
	return v1[0]*v2[0] + v1[1]*v2[1]
}

func perpendicular(v []float64) []float64 {
	return []float64{-v[1], v[0]}
}

func horizontalIntersection(v []float64, pt []float64, y float64) []float64 {
	if v[1] == 0 {
		return nil
	}
	return []float64{pt[0] + v[0]/v[1]*(y-pt[1]), y}
}

func verticalIntersection(v []float64, pt []float64, x float64) []float64 {
	if v[0] == 0 {
		return nil
	}
	return []float64{x, pt[1] + v[1]/v[0]*(x-pt[0])}
}

func closestPoint(ptA1, ptA2, ptB []float64) []float64 {
	if ptA1[0] == ptA2[0] {
		return []float64{ptA1[0], ptB[1]} // vertical vector
	}
	if ptA1[1] == ptA2[1] {
		return []float64{ptB[0], ptA1[1]} // horizontal vector
	}

	// determine which point is further away
	// we use the further point as our base in the calculation, so that the
	// vectors are more parallel, providing more accurate dot product
	v1 := []float64{ptB[0] - ptA1[0], ptB[1] - ptA1[1]}
	v2 := []float64{ptB[0] - ptA2[0], ptB[1] - ptA2[1]}
	var vFar, vA []float64
	var farPt []float64
	if dotProduct(v1, v1) > dotProduct(v2, v2) {
		vFar = v1
		vA = []float64{ptA2[0] - ptA1[0], ptA2[1] - ptA1[1]}
		farPt = ptA1
	} else {
		vFar = v2
		vA = []float64{ptA1[0] - ptA2[0], ptA1[1] - ptA2[1]}
		farPt = ptA2
	}
	// manually test if the current point can be considered to be on the line
	// If the X coordinate was on the line, would the Y coordinate be as well?
	xDist := (ptB[0] - farPt[0]) / vA[0]
	if ptB[1] == farPt[1]+xDist*vA[1] {
		return ptB
	}

	// If the Y coordinate was on the line, would the X coordinate be as well?
	yDist := (ptB[1] - farPt[1]) / vA[1]
	if ptB[0] == farPt[0]+yDist*vA[0] {
		return ptB
	}

	// current point isn't exactly on line, so return closest point
	dist := dotProduct(vA, vFar) / dotProduct(vA, vA)
	return []float64{farPt[0] + dist*vA[0], farPt[1] + dist*vA[1]}
}
