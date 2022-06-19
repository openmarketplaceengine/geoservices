package distance

import (
	"context"
	"github.com/openmarketplaceengine/geoservices"
	"time"
)

type MatrixService interface {
	GetMatrix(ctx context.Context, request PointsRequest) (*Matrix, error)
}

type PointsRequest struct {
	Origins      []geoservices.LatLng
	Destinations []geoservices.LatLng
}

type PlacesRequest struct {
	Origins      []string
	Destinations []string
}

// Matrix contains distance and duration information
// for each origin/destination pair for which a route could be calculated.
type Matrix struct {
	// OriginAddresses contains an array of addresses as returned by the API from
	// your original request.
	OriginAddresses []string `json:"origin_addresses"`
	// DestinationAddresses contains an array of addresses as returned by the API
	// from your original request.
	DestinationAddresses []string `json:"destination_addresses"`
	// Rows contains an array of elements.
	Rows []MatrixElementsRow `json:"rows"`
}

// MatrixElementsRow is a row of distance elements.
type MatrixElementsRow struct {
	Elements []MatrixElement `json:"elements"`
}

// MatrixElement is the travel distance and time for a pair of origin and
// destination.
type MatrixElement struct {
	Status string `json:"status"`
	// Duration is the length of time it takes to travel this route.
	Duration time.Duration `json:"duration"`
	// DurationInTraffic is the length of time it takes to travel this route
	// considering traffic.
	DurationInTraffic time.Duration `json:"duration_in_traffic"`
	// Distance is the total distance (in meters) of this route.
	Distance int `json:"distance"`
}
