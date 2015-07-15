// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

import (
	"strconv"
)

type Shape struct {
	Id     string
	Points ShapePoints
}

type ShapePoint struct {
	Lat           float32
	Lon           float32
	Sequence      int
	Dist_traveled float32
}

// Get a string representation of a ShapePoint
func (shape ShapePoint) String() string {
	return strconv.FormatFloat(float64(shape.Lat), 'f', 8, 32) + "," + strconv.FormatFloat(float64(shape.Lon), 'f', 8, 32)
}

// Get a string representation of this shape
func (shape Shape) String() string {
	ret := ""
	first := true
	for _, point := range shape.Points {
		if !first {
			ret += "\n"
		}
		first = false
		ret += point.String()
	}

	return ret
}

type ShapePoints []*ShapePoint

func (shapePoints ShapePoints) Len() int {
	return len(shapePoints)
}

func (shapePoints ShapePoints) Less(i, j int) bool {
	return shapePoints[i].Sequence < shapePoints[j].Sequence
}

func (shapePoints ShapePoints) Swap(i, j int) {
	shapePoints[i], shapePoints[j] = shapePoints[j], shapePoints[i]
}
