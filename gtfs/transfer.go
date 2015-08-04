// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

type Transfer struct {
	From_stop            *Stop
	To_stop              *Stop
	Transfer_type       int
	Min_transfer_time   int
}
