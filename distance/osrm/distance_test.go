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

func TestBetweenPoints(t *testing.T) {
	a := geoservices.LatLng{
		Lat: 40.791680675548136,
		Lng: -73.9650115649754,
	}
	b := geoservices.LatLng{
		Lat: 40.76866089218841,
		Lng: -73.98145413365043,
	}

	ctx := context.Background()
	res, err := BetweenPoints(ctx, new(http.Client), distance.BetweenPointsInput{
		Destinations: []geoservices.LatLng{a},
		Origins:      []geoservices.LatLng{b},
	})
	require.NoError(t, err)
	require.Len(t, res.Rows, 1)
	require.Len(t, res.Rows[0].Elements, 1)
	require.Greater(t, res.Rows[0].Elements[0].Duration, time.Duration(0))
}
