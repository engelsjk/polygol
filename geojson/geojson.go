package geojson

// The following code is a minimal implementation of GeoJSON unmarshalling for Polygons and MultiPolygons.
// It is lightly modified but mostly taken from paulmach/go.geojson; see included LICENSE file as needed.

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Geometry struct {
	Type         string `json:"type"`
	Polygon      [][][]float64
	MultiPolygon [][][][]float64
}

func (g *Geometry) UnmarshalJSON(data []byte) error {
	var object map[string]interface{}
	err := json.Unmarshal(data, &object)
	if err != nil {
		return err
	}
	t, ok := object["type"]
	if !ok {
		return errors.New("type property not defined")
	}

	if s, ok := t.(string); ok {
		g.Type = s
	} else {
		return errors.New("type property not string")
	}

	switch g.Type {
	case "Polygon":
		g.Polygon, err = decodePathSet(object["coordinates"])
	case "MultiPolygon":
		g.MultiPolygon, err = decodePolygonSet(object["coordinates"])
	}

	return err
}

func decodePosition(data interface{}) ([]float64, error) {
	coords, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid position, got %v", data)
	}

	result := make([]float64, 0, len(coords))
	for _, coord := range coords {
		if f, ok := coord.(float64); ok {
			result = append(result, f)
		} else {
			return nil, fmt.Errorf("not a valid coordinate, got %v", coord)
		}
	}

	return result, nil
}

func decodePositionSet(data interface{}) ([][]float64, error) {
	points, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid set of positions, got %v", data)
	}

	result := make([][]float64, 0, len(points))
	for _, point := range points {
		if p, err := decodePosition(point); err == nil {
			result = append(result, p)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodePathSet(data interface{}) ([][][]float64, error) {
	sets, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid path, got %v", data)
	}

	result := make([][][]float64, 0, len(sets))

	for _, set := range sets {
		if s, err := decodePositionSet(set); err == nil {
			result = append(result, s)
		} else {
			return nil, err
		}
	}

	return result, nil
}

func decodePolygonSet(data interface{}) ([][][][]float64, error) {
	polygons, ok := data.([]interface{})
	if !ok {
		return nil, fmt.Errorf("not a valid polygon, got %v", data)
	}

	result := make([][][][]float64, 0, len(polygons))
	for _, polygon := range polygons {
		if p, err := decodePathSet(polygon); err == nil {
			result = append(result, p)
		} else {
			return nil, err
		}
	}

	return result, nil
}

type Feature struct {
	Type     string    `json:"type"`
	Geometry *Geometry `json:"geometry"`
}

type FeatureCollection struct {
	Type     string     `json:"type"`
	Features []*Feature `json:"features"`
}

func UnmarshalFeatureCollection(data []byte) (*FeatureCollection, error) {
	fc := &FeatureCollection{}
	err := json.Unmarshal(data, fc)
	if err != nil {
		return nil, err
	}
	return fc, nil
}

func UnmarshalFeature(data []byte) (*Feature, error) {
	fc := &Feature{}
	err := json.Unmarshal(data, fc)
	if err != nil {
		return nil, err
	}
	return fc, nil
}
