// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

type Trip struct {
	Id                    string
	Route                 *Route
	Service               *Service
	Headsign              string
	Short_name            string
	Direction_id          int
	Block_id              string
	Shape             	  *Shape
	Wheelchair_accessible int
	Bikes_allowed         int
	StopTimes             StopTimes
	Frequencies           []*Frequency
}
