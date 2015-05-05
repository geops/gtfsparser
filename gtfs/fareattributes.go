// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

type FareAttribute struct {
	Id                string
	Price             string
	Currency_type     string
	Payment_method    int
	Transfers         int
	Transfer_duration int
	Rules             []*FareAttributeRule
}

type FareAttributeRule struct {
	Route          *Route
	Origin_id      string // connection to Zone_id in Stop
	Destination_id string // connection to Zone_id in Stop
	Contains_id    string // connection to Zone_id in Stop
}
