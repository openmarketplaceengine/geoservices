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

type Service struct {
	c *http.Client
}

func NewService(c *http.Client) *Service {
	return &Service{
		c: c,
	}
}

func (s *Service) GetMatrix(ctx context.Context, request distance.PointsRequest) (*distance.Matrix, error) {
	req := toTableRequest(request)
	res, err := table.Table(s.c, req)
	if err != nil {
		return nil, fmt.Errorf("OSRM table request error: %w", err)
	}
	return toMatrix(res), nil
}

func toTableRequest(request distance.PointsRequest) table.Request {
	var coordinates = make([]osrm.LngLat, 0)
	for _, origin := range request.Origins {
		coordinates = append(coordinates, osrm.LngLat{origin.Lng, origin.Lat})
	}
	for _, destination := range request.Destinations {
		coordinates = append(coordinates, osrm.LngLat{destination.Lng, destination.Lat})
	}
	origins := makeRange(0, len(request.Origins)-1)
	destinations := makeRange(len(request.Origins), len(request.Origins)+len(request.Destinations)-1)
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

func toMatrix(res *table.Response) *distance.Matrix {
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

	return &distance.Matrix{
		OriginAddresses:      originAddresses,
		DestinationAddresses: destinationAddresses,
		Rows:                 rows,
	}
}
