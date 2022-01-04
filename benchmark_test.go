package polygol

import (
	"fmt"
	"testing"
)

func BenchmarkAsiaUnion(b *testing.B) {
	geoms, err := loadGeoms("benchmark/asia-union/args.geojson", false)
	if err != nil {
		b.Fatal(err)
	}
	union, err := Union(Geom{}, geoms...)
	if err != nil {
		b.Fatal(err)
	}
	fmt.Printf("%+v\n", union)
}
