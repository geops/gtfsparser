// Copyright 2015 geOps
// Authors: patrick.brosi@geops.de
//
// Use of this source code is governed by a GPL v2
// license that can be found in the LICENSE file

package gtfsparser

import (
	"errors"
	"fmt"
	"github.com/geops/gtfsparser/gtfs"
	"strconv"
)

// currently, go provides no easy way to avoid repetitive
// error checking
// see http://stackoverflow.com/questions/18771569/avoid-checking-if-error-is-nil-repetition

func createAgency(r map[string]string) (*gtfs.Agency, error) {
	a := new(gtfs.Agency)
	var e error

	a.Id, e = getString("agency_id", r, false)
	if e == nil {
		a.Name, e = getString("agency_name", r, true)
	}
	if e == nil {
		a.Url, e = getString("agency_url", r, true)
	}
	if e == nil {
		a.Timezone, e = getString("agency_timezone", r, true)
	}
	if e == nil {
		a.Lang, e = getString("agency_lang", r, false)
	}
	if e == nil {
		a.Phone, e = getString("agency_phone", r, false)
	}
	if e == nil {
		a.Fare_url, e = getString("agency_fare_url", r, false)
	}

	return a, e
}

func createFeedInfo(r map[string]string) (*gtfs.FeedInfo, error) {
	f := new(gtfs.FeedInfo)
	var e error

	f.Publisher_name, e = getString("feed_publisher_name", r, true)

	if e == nil {
		f.Publisher_url, e = getString("feed_publisher_url", r, true)
	}

	if e == nil {
		f.Lang, e = getString("feed_lang", r, true)
	}

	if e == nil {
		f.Start_date, e = getString("feed_start_date", r, true)
	}

	if e == nil {
		f.End_date, e = getString("feed_end_date", r, true)
	}

	if e == nil {
		f.Version, e = getString("feed_version", r, true)
	}

	return f, e
}

func createFrequency(r map[string]string, trips *map[string]*gtfs.Trip) error {
	a := new(gtfs.Frequency)
	var e error

	var trip *gtfs.Trip

	tripid, e := getString("trip_id", r, true)

	if e != nil {
		return e
	}

	if val, ok := (*trips)[tripid]; ok {
		trip = val
	} else {
		return errors.New("No trip with id " + r["trip_id"] + " found.")
	}

	a.Exact_times, e = getBool("exact_times", r, false)

	if e == nil {
		a.Start_time, e = getString("start_time", r, true)
	}
	if e == nil {
		a.End_time, e = getString("end_time", r, true)
	}
	if e == nil {
		a.Headway_secs, e = getInt("headway_secs", r, false)
	}
	if e == nil {
		trip.Frequencies = append(trip.Frequencies, a)
	}

	return e
}

func createRoute(r map[string]string, agencies *map[string]*gtfs.Agency) (*gtfs.Route, error) {
	a := new(gtfs.Route)
	var e error

	if e == nil {
		a.Id, e = getString("route_id", r, true)
	}

	if e == nil {
		var aId, e = getString("agency_id", r, false)

		if e == nil && len(aId) != 0 {
			if val, ok := (*agencies)[aId]; ok {
				a.Agency = val
			} else {
				return nil, errors.New("No agency with id " + aId + " found.")
			}
		}
	}

	if e == nil {
		a.Short_name, e = getString("route_short_name", r, true)
	}

	if e == nil {
		a.Long_name, e = getString("route_long_name", r, true)
	}

	if e == nil {
		a.Desc, e = getString("route_desc", r, false)
	}

	if e == nil {
		a.Type, e = getInt("route_type", r, true)
	}

	if e == nil {
		a.Url, e = getString("route_url", r, false)
	}

	if e == nil {
		a.Color, e = getString("route_color", r, false)
	}

	if e == nil {
		a.Text_color, e = getString("route_text_color", r, false)
	}

	return a, e
}

func createServiceFromCalendar(r map[string]string, services *map[string]*gtfs.Service) (*gtfs.Service, error) {
	service := new(gtfs.Service)
	var e error
	service.Id, e = getString("service_id", r, true)

	// fill daybitmap
	if e == nil {
		service.Daymap[1], e = getBool("monday", r, true)
	}
	if e == nil {
		service.Daymap[2], e = getBool("tuesday", r, true)
	}
	if e == nil {
		service.Daymap[3], e = getBool("wednesday", r, true)
	}
	if e == nil {
		service.Daymap[4], e = getBool("thursday", r, true)
	}
	if e == nil {
		service.Daymap[5], e = getBool("friday", r, true)
	}
	if e == nil {
		service.Daymap[6], e = getBool("saturday", r, true)
	}
	if e == nil {
		service.Daymap[0], e = getBool("sunday", r, true)
	}
	if e == nil {
		service.Start_date, e = getDate("start_date", r, true)
	}
	if e == nil {
		service.End_date, e = getDate("end_date", r, true)
	}

	return service, e
}

func createServiceFromCalendarDates(r map[string]string, services *map[string]*gtfs.Service) (*gtfs.Service, error) {
	update := false
	var service *gtfs.Service
	var e error

	// first, check if the service already exists
	if val, ok := (*services)[r["service_id"]]; ok {
		service = val
		update = true
	} else {
		service = new(gtfs.Service)
		service.Id = r["service_id"]
	}

	// create exception
	exc := new(gtfs.ServiceException)
	var t int
	t, e = getInt("exception_type", r, true)
	exc.Type = int8(t)

	if e != nil {
		exc.Date, e = getDate("date", r, true)
	}

	service.Exceptions = append(service.Exceptions, exc)

	if update {
		return nil, e
	} else {
		return service, e
	}
}

func createStop(r map[string]string) (*gtfs.Stop, error) {
	a := new(gtfs.Stop)
	var e error

	if e == nil {
		a.Id, e = getString("stop_id", r, true)
	}

	if e == nil {
		a.Code, e = getString("stop_code", r, false)
	}

	if e == nil {
		a.Name, e = getString("stop_name", r, true)
	}

	if e == nil {
		a.Desc, e = getString("stop_desc", r, false)
	}

	if e == nil {
		a.Lat, e = getFloat("stop_lat", r, true)
	}

	if e == nil {
		a.Lon, e = getFloat("stop_lon", r, true)
	}

	if e == nil {
		a.Zone_id, e = getString("zone_id", r, false)
	}

	if e == nil {
		a.Url, e = getString("stop_url", r, false)
	}

	if e == nil {
		a.Location_type, e = getString("location_type", r, false)
	}

	if e == nil {
		a.Parent_station, e = getString("parent_station", r, false)
	}

	if e == nil {
		a.Timezone, e = getString("stop_timezone", r, false)
	}

	if e == nil {
		a.Wheelchair_boarding, e = getString("wheelchair_boarding", r, false)
	}

	return a, e
}

func createStopTime(r map[string]string, stops map[string]*gtfs.Stop, trips *map[string]*gtfs.Trip) error {
	a := new(gtfs.StopTime)

	var e error
	var trip *gtfs.Trip

	if val, ok := (*trips)[r["trip_id"]]; ok {
		trip = val
	} else {
		return errors.New("No trip with id " + r["trip_id"] + " found.")
	}

	if val, ok := stops[r["stop_id"]]; ok {
		a.Stop = val
	} else {
		return errors.New("No stop with id " + r["stop_id"] + " found.")
	}

	a.Arrival_time = r["arrival_time"]
	a.Departure_time = r["departure_time"]
	stopSeq, _ := strconv.Atoi(r["stop_sequence"])
	a.Sequence = stopSeq

	if e == nil {
		a.Headsign, e = getString("stop_headsign", r, false)
	}

	if e == nil {
		a.Pickup_type, e = getInt("pickup_type", r, false)
	}

	if e == nil {
		a.Drop_off_type, e = getInt("drop_off_type", r, false)
	}

	if e == nil {
		a.Shape_dist_traveled, e = getFloat("shape_dist_traveled", r, false)
	}

	if e == nil {
		a.Timepoint, e = getBool("Timepoint", r, false)
	}

	if e == nil {
		trip.StopTimes = append(trip.StopTimes, a)
	}

	return e
}

func createTrip(r map[string]string, routes *map[string]*gtfs.Route, services *map[string]*gtfs.Service) (*gtfs.Trip, error) {
	a := new(gtfs.Trip)
	var e error

	if e == nil {
		a.Id, e = getString("trip_id", r, true)
	}

	if rId, ok := r["route_id"]; ok {
		if val, ok := (*routes)[rId]; ok {
			a.Route = val
		} else {
			return nil, errors.New(fmt.Sprintf("No route with id %s found", rId))
		}
	}

	if sId, ok := r["service_id"]; ok {
		if val, ok := (*services)[sId]; ok {
			a.Service = val
		} else {
			return nil, errors.New(fmt.Sprintf("No service with id %s found", sId))
		}
	}

	if e == nil {
		a.Headsign, e = getString("trip_headsign", r, false)
	}

	if e == nil {
		a.Short_name, e = getString("trip_short_name", r, false)
	}

	if e == nil {
		a.Direction_id, e = getInt("direction_id", r, false)
	}

	if e == nil {
		a.Block_id, e = getString("block_id", r, false)
	}

	if e == nil {
		a.Shape_id, e = getString("shape_id", r, false)
	}

	if e == nil {
		a.Wheelchair_accessible, e = getInt("wheelchair_accessible", r, false)
	}

	if e == nil {
		a.Bikes_allowed, e = getInt("bikes_allowed", r, false)
	}

	return a, e
}

func createFareAttribute(r map[string]string) (*gtfs.FareAttribute, error) {
	a := new(gtfs.FareAttribute)
	var e error

	a.Id, e = getString("fare_id", r, true)

	if e == nil {
		a.Price, e = getString("price", r, false)
	}

	if e == nil {
		a.Currency_type, e = getString("currency_type", r, true)
	}

	if e == nil {
		a.Payment_method, e = getInt("payment_method", r, false)
	}

	if e == nil {
		a.Transfers, e = getInt("transfers", r, true)
	}

	if e == nil {
		a.Transfer_duration, e = getInt("transfer_duration", r, false)
	}

	return a, e
}

func createFareRule(r map[string]string, fareattributes *map[string]*gtfs.FareAttribute, routes *map[string]*gtfs.Route) error {
	var fareattr *gtfs.FareAttribute
	var e error
	var fareid string

	fareid, e = getString("fare_id", r, true)

	if e != nil {
		return e
	}

	// first, check if the service already exists
	if val, ok := (*fareattributes)[fareid]; ok {
		fareattr = val
	} else {
		return errors.New(fmt.Sprintf("No fare attribute with id %s found", fareid))
	}

	// create fare attribute
	rule := new(gtfs.FareAttributeRule)

	var route_id string
	route_id, e = getString("route_id", r, false)

	if e == nil && len(route_id) > 0 {
		if val, ok := (*routes)[route_id]; ok {
			rule.Route = val
		} else {
			return errors.New(fmt.Sprintf("No route with id %s found", route_id))
		}
	}

	if e != nil {
		rule.Origin_id, e = getString("origin_id", r, false)
	}

	if e != nil {
		rule.Destination_id, e = getString("destination_id", r, false)
	}

	if e != nil {
		rule.Contains_id, e = getString("contains_id", r, false)
	}

	fareattr.Rules = append(fareattr.Rules, rule)

	return nil
}

func getString(name string, r map[string]string, req bool) (string, error) {
	if val, ok := r[name]; ok {
		return val, nil
	} else if req {
		return "", errors.New(fmt.Sprintf("Expected required field %s", name))
	}
	return "", nil
}

func getInt(name string, r map[string]string, req bool) (int, error) {
	if val, ok := r[name]; ok && len(val) > 0 {
		num, err := strconv.Atoi(val)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("Expected integer for field %s, found %s", name, val))
		}
		return num, nil
	} else if req {
		return 0, errors.New(fmt.Sprintf("Expected required field %s", name))
	}
	return 0, nil
}

func getFloat(name string, r map[string]string, req bool) (float32, error) {
	if val, ok := r[name]; ok && len(val) > 0 {
		num, err := strconv.ParseFloat(val, 32)
		if err != nil {
			return 0, errors.New(fmt.Sprintf("Expected float for field %s, found %s", name, val))
		}
		return float32(num), nil
	} else if req {
		return 0, errors.New(fmt.Sprintf("Expected required field %s", name))
	}
	return 0, nil
}

func getBool(name string, r map[string]string, req bool) (bool, error) {
	if val, ok := r[name]; ok && len(val) > 0 {
		num, err := strconv.Atoi(val)
		if err != nil || (num != 0 && num != 1) {
			return false, errors.New(fmt.Sprintf("Expected 1 or 0 for field %s, found %s", name, val))
		}
		return num == 1, nil
	} else if req {
		return false, errors.New(fmt.Sprintf("Expected required field %s", name))
	}
	return false, nil
}

func getDate(name string, r map[string]string, req bool) (gtfs.Date, error) {
	var str string
	var ok bool
	if str, ok = r[name]; !ok {
		return gtfs.Date{0, 0, 0}, errors.New(fmt.Sprintf("Expected required field %s", name))
	}

	var e error
	var day, month, year int
	day, e = strconv.Atoi(str[6:8])
	if e == nil {
		month, e = strconv.Atoi(str[4:6])
	}
	if e == nil {
		year, e = strconv.Atoi(str[0:4])
	}

	if e != nil {
		return gtfs.Date{0, 0, 0}, nil
	} else {
		return gtfs.Date{int8(day), int8(month), int16(year)}, nil
	}
}
