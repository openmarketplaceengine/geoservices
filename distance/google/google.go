package google

import (
	"context"
	"github.com/openmarketplaceengine/geoservices/distance"
	"github.com/openmarketplaceengine/geoservices/geocode"
	"github.com/openmarketplaceengine/geoservices/geocode/google"
	"googlemaps.github.io/maps"
)

type Service struct {
	c *maps.Client
}

func NewService(c *maps.Client) *Service {
	return &Service{
		c: c,
	}
}

func (s *Service) GetMatrix(ctx context.Context, request distance.PointsRequest) (*distance.Matrix, error) {

	// Batch reverse-geocode all locations
	geocoder := google.NewGeocoder(s.c)
	parallelizationFactor := 10
	geocodeOut, err := geocode.BatchReverseGeocode(
		ctx,
		geocoder,
		append(request.Origins, request.Destinations...),
		parallelizationFactor)
	if err != nil {
		return nil, err
	}

	var origins []string
	var destinations []string
	for idx, e := range geocodeOut {
		if idx < len(request.Origins) {
			origins = append(origins, e.PlaceID)
		} else {
			destinations = append(destinations, e.PlaceID)
		}
	}

	res, err := BetweenPlaces(ctx, s.c, distance.PlacesRequest{
		Origins:      origins,
		Destinations: destinations,
	})
	if err != nil {
		return nil, err
	}

	return toMatrix(res, geocodeOut[:len(request.Origins)], geocodeOut[len(request.Origins):]), nil
}

func BetweenPlaces(ctx context.Context, c *maps.Client, in distance.PlacesRequest) (*maps.DistanceMatrixResponse, error) {
	var origins []string
	for _, placeID := range in.Origins {
		origins = append(origins, "place_id:"+placeID)
	}
	var destinations []string
	for _, placeID := range in.Destinations {
		destinations = append(destinations, "place_id:"+placeID)
	}
	return c.DistanceMatrix(ctx, &maps.DistanceMatrixRequest{
		Origins:                  origins,
		Destinations:             destinations,
		Mode:                     "",
		Language:                 "",
		Avoid:                    "",
		Units:                    "",
		DepartureTime:            "",
		ArrivalTime:              "",
		TrafficModel:             "",
		TransitMode:              nil,
		TransitRoutingPreference: "",
	})
}

func toMatrix(
	res *maps.DistanceMatrixResponse,
	originsOut []*geocode.ReverseGeocodeOutput,
	destinationsOut []*geocode.ReverseGeocodeOutput,
) *distance.Matrix {
	var rows []distance.MatrixElementsRow
	for i := range res.Rows {
		row := res.Rows[i]
		var elements []distance.MatrixElement
		for j := range row.Elements {
			elem := row.Elements[j]
			elements = append(elements, toElem(elem))
		}
		rows = append(rows, distance.MatrixElementsRow{Elements: elements})
	}
	var originAddresses []string
	for idx := range originsOut {
		geoResults := originsOut[idx]
		originAddresses = append(originAddresses, geoResults.FormattedAddress)
	}
	var destinationAddresses []string
	for idx := range destinationsOut {
		geoResults := destinationsOut[idx]
		destinationAddresses = append(destinationAddresses, geoResults.FormattedAddress)
	}
	return &distance.Matrix{
		OriginAddresses:      originAddresses,
		DestinationAddresses: destinationAddresses,
		Rows:                 rows,
	}
}

func toElem(res *maps.DistanceMatrixElement) distance.MatrixElement {
	return distance.MatrixElement{
		Status:            res.Status,
		Duration:          res.Duration,
		DurationInTraffic: res.DurationInTraffic,
		Distance:          res.Distance.Meters,
	}
}
