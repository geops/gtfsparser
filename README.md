# go gtfsparser

A complete*, easy to use parsing library for GTFS data. Implemented in go. Accepts folders containing GTFS files and ZIPs. Feeds are validated during parsing. ID references are transformed into pointer references where appropriate.

## Usage
    feed := gtfsparser.NewFeed()
    feed.Parse("sample-feed.zip")
    
See feed.go for exported fields.

## Example

Parsing of the [GTFS example feed](https://developers.google.com/transit/gtfs/examples/gtfs-feed):
    
```go
import (
	"github.com/geops/gtfsparser"
	"fmt"
)

func main() {
	feed := gtfsparser.NewFeed()

	feed.Parse("sample-feed.zip")

	fmt.Printf("Done, parsed %d agencies, %d stops, %d routes, %d trips, %d fare attributes\n\n",
		len(feed.Agencies), len(feed.Stops), len(feed.Routes), len(feed.Trips), len(feed.FareAttributes))

	for k, v := range feed.Stops {
        fmt.Printf("[%s] %s (@ %f,%f)\n", k, v.Name, v.Lat, v.Lon)
    }
}
```

#### Output
```
Done, parsed 1 agencies, 9 stops, 5 routes, 11 trips, 2 fare attributes

[BULLFROG] Bullfrog (Demo) (@ 36.881081,-116.817970)
[NADAV] North Ave / D Ave N (Demo) (@ 36.914894,-116.768211)
[NANAA] North Ave / N A Ave (Demo) (@ 36.914944,-116.761475)
[AMV] Amargosa Valley (Demo) (@ 36.641495,-116.400940)
[FUR_CREEK_RES] Furnace Creek Resort (Demo) (@ 36.425289,-117.133163)
[BEATTY_AIRPORT] Nye County Airport (Demo) (@ 36.868446,-116.784584)
[STAGECOACH] Stagecoach Hotel & Casino (Demo) (@ 36.915684,-116.751678)
[DADAN] Doing Ave / D Ave N (Demo) (@ 36.909489,-116.768242)
[EMSI] E Main St / S Irving St (Demo) (@ 36.905697,-116.762177)
```

## *Known restrictions

shapes.txt parsing is missing as well as parsing of transfers.txt. Validation may not be 100% complete. Tests are missing.

## License

GPL v2, see LICENSE