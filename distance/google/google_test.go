package google

import (
	"context"
	"github.com/openmarketplaceengine/geoservices"
	"github.com/openmarketplaceengine/geoservices/distance"
	"github.com/stretchr/testify/require"
	"googlemaps.github.io/maps"
	"os"
	"testing"
	"time"
)

var apiKey = os.Getenv("GOOGLE_API_KEY")

func TestGetMatrix(t *testing.T) {
	if apiKey == "" {
		t.Skip("GOOGLE_API_KEY env var is not set, skipping.")
	}

	t.Run("testGetMatrix", func(t *testing.T) {
		testGetMatrix(t)
	})

	t.Run("testGetMatrixFromPlaces", func(t *testing.T) {
		testGetMatrixFromPlaces(t)
	})

	t.Run("testGetMatrixMultipleOrigins", func(t *testing.T) {
		testGetMatrixMultipleOrigins(t)
	})
}

func testGetMatrix(t *testing.T) {
	a := geoservices.LatLng{
		Lat: 40.791680675548136,
		Lng: -73.9650115649754,
	}
	b := geoservices.LatLng{
		Lat: 40.76866089218841,
		Lng: -73.98145413365043,
	}

	ctx := context.Background()
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	require.NoError(t, err)
	matrix, err := Matrix(ctx, client, distance.MatrixPointsInput{
		Origins:      []geoservices.LatLng{a},
		Destinations: []geoservices.LatLng{b},
	})
	require.NoError(t, err)
	require.Len(t, matrix.Rows, 1)
	require.Len(t, matrix.Rows[0].Elements, 1)
	require.Greater(t, matrix.Rows[0].Elements[0].Duration, time.Duration(0))
	require.Greater(t, matrix.Rows[0].Elements[0].Distance, 0)
	require.Equal(t, matrix.OriginAddresses[0], "96 St, New York, NY 10025, USA")
	require.Equal(t, matrix.DestinationAddresses[0], "59 St-Columbus Circle, Columbus Cir, New York, NY 10023, USA")
}

func testGetMatrixFromPlaces(t *testing.T) {
	ctx := context.Background()
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	require.NoError(t, err)
	matrix, err := MatrixFromPlaces(ctx, client, distance.MatrixPlacesInput{
		Origins:      []string{"ChIJ87a7BJ5YwokR4TLbUoQMB1s"},
		Destinations: []string{"ChIJVZfjPvZYwokR-sLEBmjjniw"},
	})
	require.NoError(t, err)
	require.Len(t, matrix.Rows, 1)
	require.Len(t, matrix.Rows[0].Elements, 1)
	require.Greater(t, matrix.Rows[0].Elements[0].Duration, time.Duration(0))
	require.Greater(t, matrix.Rows[0].Elements[0].Distance, 0)
	require.Equal(t, matrix.OriginAddresses[0], "96 St, New York, NY 10025, USA")
	require.Equal(t, matrix.DestinationAddresses[0], "59 St-Columbus Circle, Columbus Cir, New York, NY 10023, USA")
}

func testGetMatrixMultipleOrigins(t *testing.T) {
	driver := geoservices.LatLng{
		Lat: 40.788088,
		Lng: -73.970988,
	}
	pickup := geoservices.LatLng{
		Lat: 40.791680675548136,
		Lng: -73.9650115649754,
	}
	destination := geoservices.LatLng{
		Lat: 40.76866089218841,
		Lng: -73.98145413365043,
	}

	ctx := context.Background()
	client, err := maps.NewClient(maps.WithAPIKey(apiKey))
	require.NoError(t, err)
	matrix, err := Matrix(ctx, client, distance.MatrixPointsInput{
		Origins:      []geoservices.LatLng{driver, pickup},
		Destinations: []geoservices.LatLng{destination},
	})
	require.NoError(t, err)
	require.Len(t, matrix.OriginAddresses, 2)
	require.Equal(t, matrix.OriginAddresses[0], "586 Columbus Ave, New York, NY 10024, USA")
	require.Equal(t, matrix.OriginAddresses[1], "96 St, New York, NY 10025, USA")
	require.Len(t, matrix.DestinationAddresses, 1)
	require.Equal(t, matrix.DestinationAddresses[0], "59 St-Columbus Circle, Columbus Cir, New York, NY 10023, USA")
	require.Len(t, matrix.Rows, 2)
	require.Len(t, matrix.Rows[0].Elements, 1)
	require.Greater(t, matrix.Rows[0].Elements[0].Duration, time.Duration(0))
	require.Greater(t, matrix.Rows[0].Elements[0].Distance, 0)

	require.Len(t, matrix.Rows[1].Elements, 1)
	require.Greater(t, matrix.Rows[1].Elements[0].Duration, time.Duration(0))
	require.Greater(t, matrix.Rows[1].Elements[0].Distance, 0)
}
