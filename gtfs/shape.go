// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

type Shape struct {
	Id     string
	Points []*ShapePoint
}

type ShapePoint struct {
	Lat           float32
	Lon           float32
	Sequence      int
	Dist_traveled float32
}
