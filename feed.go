// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfsparser

import (
	"archive/zip"
	"errors"
	"github.com/geops/gtfsparser/gtfs"
	"io"
	"os"
	opath "path"
	"sort"
)

type Feed struct {
	Agencies       map[string]*gtfs.Agency
	Stops          map[string]*gtfs.Stop
	Routes         map[string]*gtfs.Route
	Trips          map[string]*gtfs.Trip
	Services       map[string]*gtfs.Service
	FareAttributes map[string]*gtfs.FareAttribute
	Shapes         map[string]*gtfs.Shape

	zipFileCloser *zip.ReadCloser
	curFileHandle *os.File
}

// Create a new, empty feed
func NewFeed() *Feed {
	g := Feed{
		Agencies:       make(map[string]*gtfs.Agency),
		Stops:          make(map[string]*gtfs.Stop),
		Routes:         make(map[string]*gtfs.Route),
		Trips:          make(map[string]*gtfs.Trip),
		Services:       make(map[string]*gtfs.Service),
		FareAttributes: make(map[string]*gtfs.FareAttribute),
		Shapes:         make(map[string]*gtfs.Shape),
	}
	return &g
}

// Parse the GTFS data in the specified folder into the feed
func (feed *Feed) Parse(path string) error {
	var e error

	e = feed.parseAgencies(path)
	if e == nil {
		e = feed.parseStops(path)
	}
	if e == nil {
		e = feed.parseShapes(path)
	}
	if e == nil {
		e = feed.parseRoutes(path)
	}
	if e == nil {
		e = feed.parseCalendar(path)
	}
	if e == nil {
		e = feed.parseCalendarDates(path)
	}
	if e == nil {
		e = feed.parseTrips(path)
	}
	if e == nil {
		e = feed.parseStopTimes(path)
	}
	if e == nil {
		e = feed.parseFareAttributes(path)
	}
	if e == nil {
		e = feed.parseFareAttributeRules(path)
	}
	if e == nil {
		e = feed.parseFrequencies(path)
	}

	// sort stoptimes in trips
	for _, trip := range feed.Trips {
		sort.Sort(trip.StopTimes)
	}

	// sort points in shapes
	for _, shape := range feed.Shapes {
		sort.Sort(shape.Points)
	}

	// close open readers
	if feed.zipFileCloser != nil {
		feed.zipFileCloser.Close()
	}

	if feed.curFileHandle != nil {
		feed.curFileHandle.Close()
	}

	return e
}

func (feed *Feed) getFile(path string, name string) (io.Reader, error) {
	fileInfo, err := os.Stat(path)

	if err != nil {
		return nil, err
	}

	if fileInfo.IsDir() {
		if feed.curFileHandle != nil {
			// close previous handle
			feed.curFileHandle.Close()
		}

		return os.Open(opath.Join(path, name))
	} else {
		var e error
		if feed.zipFileCloser == nil {
			// reuse existing opened zip file
			feed.zipFileCloser, e = zip.OpenReader(path)
		}

		if e != nil {
			return nil, e
		}

		for _, f := range feed.zipFileCloser.File {
			if f.Name == name {
				return f.Open()
			}
		}
	}

	return nil, errors.New("Not found.")
}

func (feed *Feed) parseAgencies(path string) (err error) {
	file, e := feed.getFile(path, "agency.txt")

	if e != nil {
		return errors.New("Could not open required file agency.txt")
	}

	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"agency.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		var agency *gtfs.Agency
		agency = createAgency(record)
		feed.Agencies[agency.Id] = agency
	}

	return e
}

func (feed *Feed) parseStops(path string) (err error) {
	file, e := feed.getFile(path, "stops.txt")

	if e != nil {
		return errors.New("Could not open required file stops.txt")
	}

	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"agency.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		var stop *gtfs.Stop
		stop = createStop(record)
		if e != nil {
			return ParseError{"stops.txt", reader.Curline, e.Error()}
		}
		feed.Stops[stop.Id] = stop
	}
	return e
}

func (feed *Feed) parseRoutes(path string) (err error) {
	file, e := feed.getFile(path, "routes.txt")

	if e != nil {
		return errors.New("Could not open required file routes.txt")
	}

	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"agency.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		var route *gtfs.Route
		route = createRoute(record, feed.Agencies)
		if e != nil {
			return ParseError{"routes.txt", reader.Curline, e.Error()}
		}
		feed.Routes[route.Id] = route
	}
	return e
}

func (feed *Feed) parseCalendar(path string) (err error) {
	file, e := feed.getFile(path, "calendar.txt")

	if e != nil {
		return nil
	}

	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"agency.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		var service *gtfs.Service
		service = createServiceFromCalendar(record, feed.Services)
		if e != nil {
			return ParseError{"calendar.txt", reader.Curline, e.Error()}
		}

		// if service was parsed in-place, nil was returned
		if service != nil {
			feed.Services[service.Id] = service
		}
	}

	return e
}

func (feed *Feed) parseCalendarDates(path string) (err error) {
	file, e := feed.getFile(path, "calendar_dates.txt")

	if e != nil {
		return nil
	}

	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"calendar_dates.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		var service *gtfs.Service
		service = createServiceFromCalendarDates(record, feed.Services)

		// if service was parsed in-place, nil was returned
		if service != nil {
			feed.Services[service.Id] = service
		}
	}

	return e
}

func (feed *Feed) parseTrips(path string) (err error) {
	file, e := feed.getFile(path, "trips.txt")

	if e != nil {
		return errors.New("Could not open required file trips.txt")
	}

	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"trips.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		var trip *gtfs.Trip
		trip = createTrip(record, feed.Routes, feed.Services, feed.Shapes)
		feed.Trips[trip.Id] = trip
	}

	return e
}

func (feed *Feed) parseShapes(path string) (err error) {
	file, e := feed.getFile(path, "shapes.txt")

	if e != nil {
		return nil
	}

	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"shapes.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		createShapePoint(record, feed.Shapes)
	}

	return e
}

func (feed *Feed) parseStopTimes(path string) (err error) {
	file, e := feed.getFile(path, "stop_times.txt")

	if e != nil {
		return errors.New("Could not open required file stop_times.txt")
	}
	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"stop_times.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		createStopTime(record, feed.Stops, feed.Trips)
	}

	return e
}

func (feed *Feed) parseFrequencies(path string) (err error) {
	file, e := feed.getFile(path, "frequencies.txt")

	if e != nil {
		return nil
	}
	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"frequencies.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		createFrequency(record, feed.Trips)
	}

	return e
}

func (feed *Feed) parseFareAttributes(path string) (err error) {
	file, e := feed.getFile(path, "fare_attributes.txt")

	if e != nil {
		return nil
	}
	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"fare_attributes.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		var fa *gtfs.FareAttribute
		fa = createFareAttribute(record)
		feed.FareAttributes[fa.Id] = fa
	}

	return e
}

func (feed *Feed) parseFareAttributeRules(path string) (err error) {
	file, e := feed.getFile(path, "fare_rules.txt")

	if e != nil {
		return nil
	}
	reader := NewCsvParser(file)

	defer func() {
		if r := recover(); r != nil {
			err = ParseError{"fare_rules.txt", reader.Curline, r.(string)}
		}
	}()

	var record map[string]string
	for record = reader.ParseRecord(); record != nil; record = reader.ParseRecord() {
		createFareRule(record, feed.FareAttributes, feed.Routes)
	}

	return e
}
