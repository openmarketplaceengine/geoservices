package osrm

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/geoservices"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestReverseGeocode(t *testing.T) {
	g := NewGeocoder(new(http.Client))
	ctx := context.Background()
	out, err := g.ReverseGeocode(ctx, geoservices.LatLng{Lat: -37.813611, Lng: 144.963056})
	require.NoError(t, err)
	fmt.Println(out)
}
