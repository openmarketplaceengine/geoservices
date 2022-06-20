package google

import (
	"context"
	"github.com/openmarketplaceengine/geoservices/distance"
	"github.com/openmarketplaceengine/geoservices/geocode"
	"github.com/openmarketplaceengine/geoservices/geocode/google"
	"googlemaps.github.io/maps"
)

func GetMatrix(ctx context.Context, c *maps.Client, request distance.MatrixRequest) (*distance.Matrix, error) {

	// Batch reverse-geocode all locations
	geocoder := google.NewGeocoder(c)
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

	matrix, err := GetMatrixFromPlaces(ctx, c, distance.PlacesRequest{
		Origins:      origins,
		Destinations: destinations,
	})
	if err != nil {
		return nil, err
	}

	return matrix, nil
}

func GetMatrixFromPlaces(ctx context.Context, c *maps.Client, in distance.PlacesRequest) (*distance.Matrix, error) {
	var origins []string
	for _, placeID := range in.Origins {
		origins = append(origins, "place_id:"+placeID)
	}
	var destinations []string
	for _, placeID := range in.Destinations {
		destinations = append(destinations, "place_id:"+placeID)
	}
	matrix, err := c.DistanceMatrix(ctx, &maps.DistanceMatrixRequest{
		Origins:      origins,
		Destinations: destinations,
	})
	if err != nil {
		return nil, err
	}
	return toMatrix(matrix), err
}

func toMatrix(res *maps.DistanceMatrixResponse) *distance.Matrix {
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
	return &distance.Matrix{
		OriginAddresses:      res.OriginAddresses,
		DestinationAddresses: res.DestinationAddresses,
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
