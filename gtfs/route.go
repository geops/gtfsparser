// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

type Route struct {
	Id         string
	Agency     *Agency
	Short_name string
	Long_name  string
	Desc       string
	Type       int
	Url        string
	Color      string
	Text_color string
}