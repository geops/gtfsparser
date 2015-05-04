// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfs

import (
	"time"
)

type Service struct {
	Id         string
	Daymap     [7]bool
	Start_date Date
	End_date   Date
	Exceptions []*ServiceException
}

type ServiceException struct {
	Date Date
	Type int8
}

type Date struct {
	Day   int8
	Month int8
	Year  int16
}

func (s Service) IsActiveOn(d Date) bool {
	return (s.Daymap[int(d.GetTime().Weekday())] && !(d.GetTime().Before(s.Start_date.GetTime())) && !(d.GetTime().After(s.End_date.GetTime())) && s.GetExceptionTypeOn(d) < 2) || s.GetExceptionTypeOn(d) == 1
}

func (s Service) GetExceptionTypeOn(d Date) int8 {
	for _, e := range s.Exceptions {
		if e.Date == d {
			return e.Type
		}
	}

	return 0
}

func (d Date) GetTime() time.Time {
	return time.Date(int(d.Year), time.Month(d.Month), int(d.Day), 12, 0, 0, 0, time.UTC)
}