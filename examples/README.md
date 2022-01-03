# Examples

Both the input and output of ```polygol``` operations use a simple type ```Geom``` which represents the coordinate structure of a MultiPolygon:

```go
type Geom [][][][]float64
```

```polygol``` can be used directly with any arbitrary set of 2D coordinates, but let's consider three useful Go libraries that implement 2D geometries:

* [paulmach/go.geojson](https://github.com/paulmach/go.geojson)
* [paulmach/orb](https://github.com/paulmach/orb)
* [twpayne/go-geom](https://github.com/twpayne/go-geom)

As each of these three geometry libraries implement geometries in slightly different ways, we'll need some conversion functions to translate between their Polygon/MultiPolygon geometries and our [][][][]float64 ```Geom``` type. In the example below, the function ```g2p``` converts a library's geometry type to a ```Geom``` while ```p2g``` does the reverse. These [conversion functions](#convversion-functions) are shown in more detail below.

```go
func g2p(g *geojson.Geometry) [][][][]float64 // paulmach/go.geojson
func g2p(g orb.Geometry) [][][][]float64      // paulmach/orb
func g2p(g geom.T) [][][][]float64            // twpayne/go-geom
```

And finally, for illustrative purposes the subject geometry (A) in the example below is a single GeoJSON Feature and the clipping geometries (B) are contained in a GeoJSON FeatureCollection.

```go

var A, B [][][][]float64

// A

rawFeatureJSON := []byte(`
    { "type": "Feature",
        "geometry": {"type": "Polygon", "coordinates": [[[...]]]},
    }`)

f, _ := geojson.UnmarshalFeature(rawJSONFeature)

A = g2p(f.Geometry)

// B

rawFeatureCollectionJSON := []byte(`
    { "type": "FeatureCollection",
        "features": [
            "geometry": {"type": "Polygon", "coordinates": [[[...]]]},
            "geometry": {"type": "MultiPolygon", "coordinates": [[[[...]]]]}
        ]
    }`)

fc, _ := geojson.UnmarshalFeatureCollection(rawFeatureCollectionJSON)

B = make([]polygol.Geom, len(fc.Features))

for i := range fc.Features {
    B[i] = g2p(fc.Features[i].Geometry)
}

// C

unionC, _ := polygol.Union(A, B...)
intersectionC, _ := polygol.Intersection(A, B...)
differenceC, _ := polygol.Difference(A, B...)
xorC, _ := polygol.XOR(A, B...)

geometryUnion := p2g(unionC)
geometryIntersection := p2g(intersectionC)
geometryDifference := p2g(differenceC)
geometryXOR := p2g(xorCxorC)

```

Note: The GeoJSON unmarshalling in the example above assumes the use of either [paulmach/go.geojson](https://github.com/paulmach/go.geojson) or [paulmach/orb/geojson](https://github.com/paulmach/orb/tree/master/geojson), although using the GeoJSON decoding in [twpayne/go-geom](https://github.com/twpayne/go-geom) would not look much different.

## Conversion Functions

The conversion functions shown here are for illustrative purposes only, there are probably much more efficient ways of doing these. Your mileage may vary.

### paulmach/go.geojson

```go
func g2p(g *geojson.Geometry) [][][][]float64 {
	switch g.Type {
	case geojson.GeometryPolygon:
		return [][][][]float64{g.Polygon}
	case geojson.GeometryMultiPolygon:
		return g.MultiPolygon
	}
	return nil
}

func p2g(p [][][][]float64) *geojson.Geometry {
	return geojson.NewMultiPolygonGeometry(p...)
}
```

### paulmach/orb

```go
func g2p(g orb.Geometry) [][][][]float64 {

	var p [][][][]float64

	switch v := g.(type) {
	case orb.Polygon:
		p = make([][][][]float64, 1)
		p[0] = make([][][]float64, len(v))
		for i := range v { // rings
			p[0][i] = make([][]float64, len(v[i]))
			for j := range v[i] { // points
				pt := v[i][j]
				p[0][i][j] = []float64{pt.X(), pt.Y()}
			}
		}
	case orb.MultiPolygon:
		p = make([][][][]float64, len(v))
		for i := range v { // polygons
			p[i] = make([][][]float64, len(v[i]))
			for j := range v[i] { // rings
				p[i][j] = make([][]float64, len(v[i][j]))
				for k := range v[i][j] { // points
					pt := v[i][j][k]
					p[i][j][k] = []float64{pt.X(), pt.Y()}
				}
			}
		}
	}

	return p
}

func p2g(p [][][][]float64) orb.Geometry {

	g := make(orb.MultiPolygon, len(p))

	for i := range p {
		g[i] = make([]orb.Ring, len(p[i]))
		for j := range p[i] {
			g[i][j] = make([]orb.Point, len(p[i][j]))
			for k := range p[i][j] {
				pt := p[i][j][k]
				point := orb.Point{pt[0], pt[1]}
				g[i][j][k] = point
			}
		}
	}
	return g
}
```

### twpayne/gogeom

```go
func g2p(g geom.T) [][][][]float64 {

	var coords [][][]geom.Coord

	switch v := g.(type) {
	case *geom.Polygon:
		coords = [][][]geom.Coord{v.Coords()}
	case *geom.MultiPolygon:
		coords = v.Coords()
	}

	p := make([][][][]float64, len(coords))

	for i := range coords {
		p[i] = make([][][]float64, len(coords[i]))
		for j := range coords[i] {
			p[i][j] = make([][]float64, len(coords[i][j]))
			for k := range coords[i][j] {
				coord := coords[i][j][k]
				pt := make([]float64, 2)
				pt[0], pt[1] = coord.X(), coord.Y()
				p[i][j][k] = pt
			}
		}
	}

	return p
}

func p2g(p [][][][]float64) geom.T {

	coords := make([][][]geom.Coord, len(p))

	for i := range p {
		coords[i] = make([][]geom.Coord, len(p[i]))
		for j := range p[i] {
			coords[i][j] = make([]geom.Coord, len(p[i][j]))
			for k := range p[i][j] {
				pt := p[i][j][k]
				coord := geom.Coord{pt[0], pt[1]}
				coords[i][j][k] = coord
			}
		}
	}

	return geom.NewMultiPolygon(geom.XY).MustSetCoords(coords)
}
```
