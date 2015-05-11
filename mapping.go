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

func createAgency(r map[string]string) *gtfs.Agency {
	a := new(gtfs.Agency)

	a.Id = getString("agency_id", r, false)
	a.Name = getString("agency_name", r, true)
	a.Url = getString("agency_url", r, true)
	a.Timezone = getString("agency_timezone", r, true)
	a.Lang = getString("agency_lang", r, false)
	a.Phone = getString("agency_phone", r, false)
	a.Fare_url = getString("agency_fare_url", r, false)

	return a
}

func createFeedInfo(r map[string]string) *gtfs.FeedInfo {
	f := new(gtfs.FeedInfo)

	f.Publisher_name = getString("feed_publisher_name", r, true)
	f.Publisher_url = getString("feed_publisher_url", r, true)
	f.Lang = getString("feed_lang", r, true)
	f.Start_date = getString("feed_start_date", r, true)
	f.End_date = getString("feed_end_date", r, true)
	f.Version = getString("feed_version", r, true)

	return f
}

func createFrequency(r map[string]string, trips map[string]*gtfs.Trip) {
	a := new(gtfs.Frequency)
	var trip *gtfs.Trip

	tripid := getString("trip_id", r, true)

	if val, ok := trips[tripid]; ok {
		trip = val
	} else {
		panic("No trip with id " + r["trip_id"] + " found.")
	}

	a.Exact_times = getBool("exact_times", r, false)
	a.Start_time = getString("start_time", r, true)
	a.End_time = getString("end_time", r, true)
	a.Headway_secs = getInt("headway_secs", r, false)
	trip.Frequencies = append(trip.Frequencies, a)
}

func createRoute(r map[string]string, agencies map[string]*gtfs.Agency) *gtfs.Route {
	a := new(gtfs.Route)
	a.Id = getString("route_id", r, true)

	var aId = getString("agency_id", r, false)

	if len(aId) != 0 {
		if val, ok := agencies[aId]; ok {
			a.Agency = val
		} else {
			panic("No agency with id " + aId + " found.")
		}
	}

	a.Short_name = getString("route_short_name", r, true)
	a.Long_name = getString("route_long_name", r, true)
	a.Desc = getString("route_desc", r, false)
	a.Type = getInt("route_type", r, true)
	a.Url = getString("route_url", r, false)
	a.Color = getString("route_color", r, false)
	a.Text_color = getString("route_text_color", r, false)

	return a
}

func createServiceFromCalendar(r map[string]string, services map[string]*gtfs.Service) *gtfs.Service {
	service := new(gtfs.Service)
	service.Id = getString("service_id", r, true)

	// fill daybitmap
	service.Daymap[1] = getBool("monday", r, true)
	service.Daymap[2] = getBool("tuesday", r, true)
	service.Daymap[3] = getBool("wednesday", r, true)
	service.Daymap[4] = getBool("thursday", r, true)
	service.Daymap[5] = getBool("friday", r, true)
	service.Daymap[6] = getBool("saturday", r, true)
	service.Daymap[0] = getBool("sunday", r, true)
	service.Start_date = getDate("start_date", r, true)
	service.End_date = getDate("end_date", r, true)

	return service
}

func createServiceFromCalendarDates(r map[string]string, services map[string]*gtfs.Service) *gtfs.Service {
	update := false
	var service *gtfs.Service

	// first, check if the service already exists
	if val, ok := services[r["service_id"]]; ok {
		service = val
		update = true
	} else {
		service = new(gtfs.Service)
		service.Id = r["service_id"]
	}

	// create exception
	exc := new(gtfs.ServiceException)
	var t int
	t = getInt("exception_type", r, true)
	exc.Type = int8(t)
	exc.Date = getDate("date", r, true)

	service.Exceptions = append(service.Exceptions, exc)

	if update {
		return nil
	} else {
		return service
	}
}

func createStop(r map[string]string) *gtfs.Stop {
	a := new(gtfs.Stop)

	a.Id = getString("stop_id", r, true)
	a.Code = getString("stop_code", r, false)
	a.Name = getString("stop_name", r, true)
	a.Desc = getString("stop_desc", r, false)
	a.Lat = getFloat("stop_lat", r, true)
	a.Lon = getFloat("stop_lon", r, true)
	a.Zone_id = getString("zone_id", r, false)
	a.Url = getString("stop_url", r, false)
	a.Location_type = getString("location_type", r, false)
	a.Parent_station = getString("parent_station", r, false)
	a.Timezone = getString("stop_timezone", r, false)
	a.Wheelchair_boarding = getString("wheelchair_boarding", r, false)

	return a
}

func createStopTime(r map[string]string, stops map[string]*gtfs.Stop, trips map[string]*gtfs.Trip) {
	a := new(gtfs.StopTime)
	var trip *gtfs.Trip

	if val, ok := trips[getString("trip_id", r, true)]; ok {
		trip = val
	} else {
		panic("No trip with id " + getString("trip_id", r, true) + " found.")
	}

	if val, ok := stops[getString("stop_id", r, true)]; ok {
		a.Stop = val
	} else {
		panic("No stop with id " + getString("stop_id", r, true) + " found.")
	}

	a.Arrival_time = getString("arrival_time", r, true)
	a.Departure_time = getString("departure_time", r, true)
	a.Sequence = getInt("stop_sequence", r, true)
	a.Headsign = getString("stop_headsign", r, false)
	a.Pickup_type = getInt("pickup_type", r, false)
	a.Drop_off_type = getInt("drop_off_type", r, false)
	a.Shape_dist_traveled = getFloat("shape_dist_traveled", r, false)
	a.Timepoint = getBool("Timepoint", r, false)

	trip.StopTimes = append(trip.StopTimes, a)

}

func createTrip(r map[string]string, routes map[string]*gtfs.Route, 
	services map[string]*gtfs.Service,
	shapes map[string]*gtfs.Shape) *gtfs.Trip {
	a := new(gtfs.Trip)
	a.Id = getString("trip_id", r, true)

	if rId, ok := r["route_id"]; ok {
		if val, ok := routes[rId]; ok {
			a.Route = val
		} else {
			panic(fmt.Sprintf("No route with id %s found", rId))
		}
	}

	if sId, ok := r["service_id"]; ok {
		if val, ok := services[sId]; ok {
			a.Service = val
		} else {
			panic(fmt.Sprintf("No service with id %s found", sId))
		}
	}

	a.Headsign = getString("trip_headsign", r, false)
	a.Short_name = getString("trip_short_name", r, false)
	a.Direction_id = getInt("direction_id", r, false)
	a.Block_id = getString("block_id", r, false)

	shapeId := getString("shape_id", r, false)

    if len(shapeId) > 0 {
		if val, ok := shapes[shapeId]; ok {
			a.Shape = val
		} else {
			panic(fmt.Sprintf("No shape with id %s found", shapeId))
		}
	}

	a.Wheelchair_accessible = getInt("wheelchair_accessible", r, false)
	a.Bikes_allowed = getInt("bikes_allowed", r, false)

	return a
}

func createShapePoint(r map[string]string, shapes map[string]*gtfs.Shape) {
	shapeId := getString("shape_id", r, true)
	var shape *gtfs.Shape

	if val, ok := shapes[shapeId]; ok {
		shape = val
	} else {
		// create new shape
		shape = new(gtfs.Shape)
		// push it onto the shape map
		shapes[shapeId] = shape
	}

	shape.Points = append(shape.Points, &gtfs.ShapePoint{
		Lat: getFloat("shape_pt_lat", r, true),
		Lon: getFloat("shape_pt_lon", r, true),
		Sequence: getInt("shape_pt_sequence", r, true),
		Dist_traveled: getFloat("shape_dist_traveled", r, false),
	})
}

func createFareAttribute(r map[string]string) *gtfs.FareAttribute {
	a := new(gtfs.FareAttribute)

	a.Id = getString("fare_id", r, true)
	a.Price = getString("price", r, false)
	a.Currency_type = getString("currency_type", r, true)
	a.Payment_method = getInt("payment_method", r, false)
	a.Transfers = getInt("transfers", r, true)
	a.Transfer_duration = getInt("transfer_duration", r, false)

	return a
}

func createFareRule(r map[string]string, fareattributes map[string]*gtfs.FareAttribute, routes map[string]*gtfs.Route) {
	var fareattr *gtfs.FareAttribute
	var fareid string

	fareid = getString("fare_id", r, true)

	// first, check if the service already exists
	if val, ok := fareattributes[fareid]; ok {
		fareattr = val
	} else {
		panic(fmt.Sprintf("No fare attribute with id %s found", fareid))
	}

	// create fare attribute
	rule := new(gtfs.FareAttributeRule)

	var route_id string
	route_id = getString("route_id", r, false)

	if len(route_id) > 0 {
		if val, ok := routes[route_id]; ok {
			rule.Route = val
		} else {
			panic(fmt.Sprintf("No route with id %s found", route_id))
		}
	}

	rule.Origin_id = getString("origin_id", r, false)
	rule.Destination_id = getString("destination_id", r, false)
	rule.Contains_id = getString("contains_id", r, false)

	fareattr.Rules = append(fareattr.Rules, rule)
}

func getString(name string, r map[string]string, req bool) string {
	if val, ok := r[name]; ok {
		return val
	} else if req {
		panic(fmt.Sprintf("Expected required field %s", name))
	}
	return ""
}

func getInt(name string, r map[string]string, req bool) int {
	if val, ok := r[name]; ok && len(val) > 0 {
		num, err := strconv.Atoi(val)
		if err != nil {
			panic(errors.New(fmt.Sprintf("Expected integer for field %s, found %s", name, val)))
		}
		return num
	} else if req {
		panic(fmt.Sprintf("Expected required field %s", name))
	}
	return 0
}

func getFloat(name string, r map[string]string, req bool) float32 {
	if val, ok := r[name]; ok && len(val) > 0 {
		num, err := strconv.ParseFloat(val, 32)
		if err != nil {
			panic(fmt.Sprintf("Expected float for field %s, found %s", name, val))
		}
		return float32(num)
	} else if req {
		panic(fmt.Sprintf("Expected required field %s", name))
	}
	return 0
}

func getBool(name string, r map[string]string, req bool) bool {
	if val, ok := r[name]; ok && len(val) > 0 {
		num, err := strconv.Atoi(val)
		if err != nil || (num != 0 && num != 1) {
			panic(fmt.Sprintf("Expected 1 or 0 for field %s, found %s", name, val))
		}
		return num == 1
	} else if req {
		panic(fmt.Sprintf("Expected required field %s", name))
	}
	return false
}

func getDate(name string, r map[string]string, req bool) gtfs.Date {
	var str string
	var ok bool
	if str, ok = r[name]; !ok {
		panic(fmt.Sprintf("Expected required field %s", name))
	}

	var day, month, year int
	var e error
	day, e = strconv.Atoi(str[6:8])
	if e == nil {
		month, e = strconv.Atoi(str[4:6])
	}
	if e == nil {
		year, e = strconv.Atoi(str[0:4])
	}

	if e != nil {
		panic(e)
	} else {
		return gtfs.Date{int8(day), int8(month), int16(year)}
	}
}
