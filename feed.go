// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfsparser

import (
	"fmt"
	"github.com/geops/gtfsparser/gtfs"
	"sort"
)

type Feed struct {
	Agencies map[string]*gtfs.Agency
	Stops    map[string]*gtfs.Stop
	Routes   map[string]*gtfs.Route
	Trips    map[string]*gtfs.Trip
	Services map[string]*gtfs.Service
}

func NewFeed() *Feed {
	g := Feed{
		Agencies: make(map[string]*gtfs.Agency),
		Stops:    make(map[string]*gtfs.Stop),
		Routes:   make(map[string]*gtfs.Route),
		Trips:    make(map[string]*gtfs.Trip),
		Services: make(map[string]*gtfs.Service),
	}
	return &g
}

func (feed *Feed) Parse(folder string) error {
	fmt.Println("Parsing GTFS at " + folder)

	var e error

	e = feed.ParseAgencies(folder + "/agency.txt")

	if e == nil {
		e = feed.ParseStops(folder + "/stops.txt")
	}
	if e == nil {
		e = feed.ParseRoutes(folder + "/routes.txt")
	}
	if e == nil {
		e = feed.ParseCalendar(folder + "/calendar.txt")
	}
	if e == nil {
		e = feed.ParseCalendarDates(folder + "/calendar_dates.txt")
	}
	if e == nil {
		e = feed.ParseTrips(folder + "/trips.txt")
	}
	if e == nil {
		e = feed.ParseStopTimes(folder + "/stop_times.txt")
	}
	if e == nil {
		e = feed.ParseFrequencies(folder + "/frequencies.txt")
	} else {
		return e
	}

	for _, trip := range feed.Trips {
		sort.Sort(trip.StopTimes)
	}

	fmt.Printf("Done, parse %d agencies, %d stops, %d routes, %d trips\n",
		len(feed.Agencies), len(feed.Stops), len(feed.Routes), len(feed.Trips))

	return nil
}

func (feed *Feed) ParseAgencies(file string) error {
	reader, e := NewCsvParser(file)

	if e != nil {
		return e
	}

	var record map[string]string
	for record, e = reader.ParseRecord(); record != nil; record, e = reader.ParseRecord() {
		var agency *gtfs.Agency
		agency, e = CreateAgency(record)
		if e != nil {
			return ParseError{file, reader.Curline, e.Error()}
		}
		feed.Agencies[agency.Id] = agency
	}
	return e
}

func (feed *Feed) ParseStops(file string) error {
	reader, e := NewCsvParser(file)

	if e != nil {
		return e
	}

	var record map[string]string
	for record, e = reader.ParseRecord(); record != nil; record, e = reader.ParseRecord() {
		var stop *gtfs.Stop
		stop, e = CreateStop(record)
		if e != nil {
			return ParseError{file, reader.Curline, e.Error()}
		}
		feed.Stops[stop.Id] = stop
	}
	return e
}

func (feed *Feed) ParseRoutes(file string) error {
	reader, e := NewCsvParser(file)

	if e != nil {
		return e
	}

	var record map[string]string
	for record, e = reader.ParseRecord(); record != nil; record, e = reader.ParseRecord() {
		var route *gtfs.Route
		route, e = CreateRoute(record, &feed.Agencies)
		if e != nil {
			return ParseError{file, reader.Curline, e.Error()}
		}
		feed.Routes[route.Id] = route
	}
	return e
}

func (feed *Feed) ParseCalendar(file string) error {
	reader, e := NewCsvParser(file)

	if e != nil {
		// we dont require calendar.txt, there are many feeds that
		// dont use it an entirely rely on calendar_dates.txt
		return nil
	}

	var record map[string]string
	for record, e = reader.ParseRecord(); record != nil; record, e = reader.ParseRecord() {
		var service *gtfs.Service
		service, e = CreateServiceFromCalendar(record, &feed.Services)
		if e != nil {
			return ParseError{file, reader.Curline, e.Error()}
		}

		// if service was parsed in-place, nil was returned
		if service != nil {
			feed.Services[service.Id] = service
		}
	}

	return e
}

func (feed *Feed) ParseCalendarDates(file string) error {
	reader, e := NewCsvParser(file)

	if e != nil {
		return nil
	}

	var record map[string]string
	for record, e = reader.ParseRecord(); record != nil; record, e = reader.ParseRecord() {
		var service *gtfs.Service
		service, e = CreateServiceFromCalendarDates(record, &feed.Services)

		if e != nil {
			return ParseError{file, reader.Curline, e.Error()}
		}

		// if service was parsed in-place, nil was returned
		if service != nil {
			feed.Services[service.Id] = service
		}
	}

	return e
}

func (feed *Feed) ParseTrips(file string) error {
	reader, e := NewCsvParser(file)

	if e != nil {
		return e
	}

	var record map[string]string
	for record, e = reader.ParseRecord(); record != nil; record, e = reader.ParseRecord() {
		var trip *gtfs.Trip
		trip, e = CreateTrip(record, &feed.Routes, &feed.Services)
		if e != nil {
			return ParseError{file, reader.Curline, e.Error()}
		}
		feed.Trips[trip.Id] = trip
	}

	return e
}

func (feed *Feed) ParseStopTimes(file string) error {
	reader, e := NewCsvParser(file)

	if e != nil {
		return e
	}

	var record map[string]string
	for record, e = reader.ParseRecord(); record != nil; record, e = reader.ParseRecord() {
		e = CreateStopTime(record, feed.Stops, &feed.Trips)
		if e != nil {
			return ParseError{file, reader.Curline, e.Error()}
		}
	}

	return e
}

func (feed *Feed) ParseFrequencies(file string) error {
	reader, e := NewCsvParser(file)

	if e == nil {
		var record map[string]string
		for record, e = reader.ParseRecord(); record != nil; record, e = reader.ParseRecord() {
			e = CreateFrequency(record, &feed.Trips)
			if e != nil {
				return ParseError{file, reader.Curline, e.Error()}
			}
		}
	}

	return e
}
