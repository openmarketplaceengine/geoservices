package osrm

import (
	"context"
	"github.com/openmarketplaceengine/geoservices"
	"github.com/openmarketplaceengine/geoservices/distance"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestGetMatrix(t *testing.T) {
	a := geoservices.LatLng{
		Lat: 40.791680675548136,
		Lng: -73.9650115649754,
	}
	b := geoservices.LatLng{
		Lat: 40.76866089218841,
		Lng: -73.98145413365043,
	}

	ctx := context.Background()
	matrix, err := GetMatrix(ctx, &http.Client{}, distance.PointsRequest{
		Origins:      []geoservices.LatLng{a},
		Destinations: []geoservices.LatLng{b},
	})
	require.NoError(t, err)
	require.Len(t, matrix.Rows, 1)
	require.Len(t, matrix.Rows[0].Elements, 1)
	require.Greater(t, matrix.Rows[0].Elements[0].Duration, time.Duration(0))
	require.Greater(t, matrix.Rows[0].Elements[0].Distance, 0)
}
