package polygol

import (
	"testing"
)

func BenchmarkAsiaUnion(b *testing.B) {
	geoms, err := loadGeoms("test/end-to-end/countries-asia/args.geojson")
	if err != nil {
		b.Fatal(err)
	}
	_, err = Union(Geom{}, geoms...)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkAfricaUnion(b *testing.B) {
	geoms, err := loadGeoms("test/end-to-end/countries-africa/args.geojson")
	if err != nil {
		b.Fatal(err)
	}
	_, err = Union(Geom{}, geoms...)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkEuropeUnion(b *testing.B) {
	geoms, err := loadGeoms("test/end-to-end/countries-europe/args.geojson")
	if err != nil {
		b.Fatal(err)
	}
	_, err = Union(Geom{}, geoms...)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkNorthAmericaUnion(b *testing.B) {
	geoms, err := loadGeoms("test/end-to-end/countries-north-america/args.geojson")
	if err != nil {
		b.Fatal(err)
	}
	_, err = Union(Geom{}, geoms...)
	if err != nil {
		b.Fatal(err)
	}
}

func BenchmarkSouthAmericaUnion(b *testing.B) {
	geoms, err := loadGeoms("test/end-to-end/countries-south-america/args.geojson")
	if err != nil {
		b.Fatal(err)
	}
	_, err = Union(Geom{}, geoms...)
	if err != nil {
		b.Fatal(err)
	}
}
