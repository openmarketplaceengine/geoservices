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

func TestGetMatrix(t *testing.T) {
	apiKey := os.Getenv("GOOGLE_API_KEY")
	require.NotEmpty(t, apiKey)

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
	s := NewService(client)
	matrix, err := s.GetMatrix(ctx, distance.PointsRequest{
		Origins:      []geoservices.LatLng{a},
		Destinations: []geoservices.LatLng{b},
	})
	require.NoError(t, err)
	require.Len(t, matrix.Rows, 1)
	require.Len(t, matrix.Rows[0].Elements, 1)
	require.Greater(t, matrix.Rows[0].Elements[0].Duration, time.Duration(0))
	require.Greater(t, matrix.Rows[0].Elements[0].Distance, 0)
}
