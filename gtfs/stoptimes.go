// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

type StopTime struct {
	Arrival_time        string
	Departure_time      string
	Stop                *Stop
	Sequence            int
	Headsign            string
	Pickup_type         int
	Drop_off_type       int
	Shape_dist_traveled float32
	Timepoint           bool
}

type StopTimes []*StopTime

func (stopTimes StopTimes) Len() int {
	return len(stopTimes)
}

func (stopTimes StopTimes) Less(i, j int) bool {
	return stopTimes[i].Sequence < stopTimes[j].Sequence
}

func (stopTimes StopTimes) Swap(i, j int) {
	stopTimes[i], stopTimes[j] = stopTimes[j], stopTimes[i]
}
