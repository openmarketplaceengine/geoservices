package osrm

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/geoservices/distance"
	"github.com/openmarketplaceengine/geoservices/osrm"
	"github.com/openmarketplaceengine/geoservices/osrm/table"
	"net/http"
	"time"
)

func BetweenPoints(ctx context.Context, c *http.Client, in distance.BetweenPointsInput) (*distance.MatrixResponse, error) {
	request := toTableRequest(in)
	res, err := table.Table(c, request)
	if err != nil {
		return nil, fmt.Errorf("OSRM table request error: %w", err)
	}
	return toMatrixResponse(res), nil
}

func toTableRequest(in distance.BetweenPointsInput) table.Request {
	var coordinates = make([]osrm.LngLat, 0)
	for _, origin := range in.Origins {
		coordinates = append(coordinates, osrm.LngLat{origin.Lng, origin.Lat})
	}
	for _, destination := range in.Destinations {
		coordinates = append(coordinates, osrm.LngLat{destination.Lng, destination.Lat})
	}
	origins := makeRange(0, len(in.Origins)-1)
	destinations := makeRange(len(in.Origins), len(in.Origins)+len(in.Destinations)-1)
	return table.Request{
		Coordinates:  coordinates,
		Origins:      origins,
		Destinations: destinations,
		Annotations:  table.DurationDistance,
	}
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func toMatrixResponse(res *table.Response) *distance.MatrixResponse {
	var rows []distance.MatrixElementsRow

	var originAddresses []string
	var destinationAddresses []string
	for i, source := range res.Sources {
		originAddresses = append(originAddresses, source.Location.Textual())
		var elements []distance.MatrixElement
		for j, destination := range res.Destinations {
			destinationAddresses = append(destinationAddresses, destination.Location.Textual())
			elements = append(elements, distance.MatrixElement{
				Status:            "",
				Duration:          time.Duration(res.Durations[i][j]) * time.Second,
				DurationInTraffic: 0,
				Distance:          int(res.Distances[i][j] * 1000),
			})
		}
		rows = append(rows, distance.MatrixElementsRow{Elements: elements})
	}

	return &distance.MatrixResponse{
		OriginAddresses:      originAddresses,
		DestinationAddresses: destinationAddresses,
		Rows:                 rows,
	}
}
