// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

type Stop struct {
	Id                  string
	Code                string
	Name                string
	Desc                string
	Lat                 float32
	Lon                 float32
	Zone_id             string
	Url                 string
	Location_type       int
	Parent_station      string
	Timezone            string
	Wheelchair_boarding int
}
