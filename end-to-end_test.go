package polygol

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"testing"

	"github.com/engelsjk/polygol/geojson"
)

const (
	endToEndDir = "test/end-to-end"
)

var (
	// USE ME TO RUN ONLY ONE TEST
	targetOnly = ""
	opOnly     = ""

	// USE ME TO SKIP TESTS
	targetsSkip = []string{}
	opsSkip     = []string{}
)

type TestCase struct {
	Name          string
	OperationType string
	ResultPath    string
}

func TestEndToEnd(t *testing.T) {

	targets, err := ioutil.ReadDir(endToEndDir)
	if err != nil {
		t.Fatal(err)
	}

	for _, target := range targets {

		if contains(targetsSkip, target.Name()) {
			fmt.Printf("skipping target %s...\n", target.Name())
			continue
		}

		if !target.IsDir() {
			continue
		}

		targetDir := path.Join(endToEndDir, target.Name())

		argsPath := path.Join(targetDir, "args.geojson")

		args, err := loadGeoms(argsPath, false)
		if err != nil {
			t.Fatal(err)
		}

		files, err := ioutil.ReadDir(targetDir)
		if err != nil {
			log.Fatal(err)
		}

		testCases := []TestCase{}

		for _, f := range files {
			if f.Name() == "args.geojson" {
				continue
			}
			ext := filepath.Ext(f.Name())
			if ext != ".geojson" {
				continue
			}
			fn := f.Name()
			opType := strings.TrimSuffix(fn, ext)
			fp := filepath.Join(targetDir, fn)
			if opType != "all" {
				testCases = append(testCases, TestCase{
					Name:          fmt.Sprintf("%s-all", target.Name()),
					OperationType: opType,
					ResultPath:    fp,
				})
			} else {
				testCases = []TestCase{
					{
						Name:          fmt.Sprintf("%s-union", target.Name()),
						OperationType: "union",
						ResultPath:    fp,
					},
					{
						Name:          fmt.Sprintf("%s-intersection", target.Name()),
						OperationType: "intersection",
						ResultPath:    fp,
					},
					{
						Name:          fmt.Sprintf("%s-all", target.Name()),
						OperationType: "xor",
						ResultPath:    fp,
					},
					{
						Name:          fmt.Sprintf("%s-difference", target.Name()),
						OperationType: "difference",
						ResultPath:    fp,
					},
				}
			}
		}

		for _, testCase := range testCases {

			t.Run(testCase.Name, func(t *testing.T) {

				t.Parallel() // run all end-to-end tests in parallel

				if contains(opsSkip, testCase.OperationType) {
					fmt.Printf("skipping op type %s...\n", testCase.OperationType)
				}

				geoms, err := loadGeoms(testCase.ResultPath, true)
				if err != nil {
					t.Fatal(err)
				}

				expected := geoms[0]

				result, err := newOperation(testCase.OperationType).run(args[0], args[1:]...)
				if err != nil {
					t.Error(err)
				}

				fmt.Printf("%+v\n", result)

				expect(t, equalMultiPoly(expected, result))
			})
		}
	}
}

func TestAsiaUnion(t *testing.T) {

	geoms, err := loadGeoms("test/end-to-end/asia-union/args.geojson", false)
	if err != nil {
		log.Fatal(err)
	}

	union, err := Union(Geom{}, geoms...)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", union)
}

func contains(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func loadGeoms(filepath string, singleFeature bool) ([]Geom, error) {

	f, err := os.Open(filepath)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}

	var features []*geojson.Feature

	if singleFeature {
		f, err := geojson.UnmarshalFeature(b)
		if err != nil {
			return nil, err
		}
		features = append(features, f)
	} else {
		fc, err := geojson.UnmarshalFeatureCollection(b)
		if err != nil {
			return nil, err
		}
		features = append(features, fc.Features...)
	}

	geoms := make([]Geom, len(features))
	for i := range features {
		fg := features[i].Geometry
		switch fg.Type {
		case "Polygon":
			geoms[i] = Geom{fg.Polygon}
		case "MultiPolygon":
			geoms[i] = fg.MultiPolygon
		default:
			return nil, fmt.Errorf("only polygon or multipolygon geometry types supported")
		}
	}

	return geoms, nil
}
