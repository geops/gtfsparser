// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

type Agency struct {
	Id       string
	Name     string
	Url      string
	Timezone string
	Lang     string
	Phone    string
	Fare_url string
}
