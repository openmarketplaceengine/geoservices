package ranking

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/geoservices"
	"github.com/stretchr/testify/require"
	"math/rand"
	"testing"
)

func randPoint(sw, ne geoservices.LatLng) geoservices.LatLng {
	return geoservices.LatLng{
		Lat: sw.Lat + (ne.Lat-sw.Lat)*rand.Float64(),
		Lng: sw.Lng + (ne.Lng-sw.Lng)*rand.Float64(),
	}
}

func TestRank(t *testing.T) {
	tests := map[string]struct {
		buildOrigins func() map[geoservices.LatLng]string
		destination  geoservices.LatLng
		expect       func(t *testing.T, out []geoservices.LatLng, originsMap map[geoservices.LatLng]string)
	}{
		"Load test": {
			buildOrigins: func() map[geoservices.LatLng]string {
				n := 100
				out := make(map[geoservices.LatLng]string, n)
				southWest := geoservices.LatLng{40.67374692754084, -74.01655691637214}
				northEast := geoservices.LatLng{40.74272240069717, -73.93601257175355}
				for i := 0; i < n; i++ {
					rp := randPoint(southWest, northEast)
					out[rp] = fmt.Sprintf("%d", i)
				}
				return out
			},
			// Brooklyn City Dental
			destination: geoservices.LatLng{40.7263248173875, -73.95246643844668},
			expect: func(t *testing.T, out []geoservices.LatLng, originsMap map[geoservices.LatLng]string) {
				require.NotEmpty(t, out)
			},
		},
		"Brooklyn City Dental": {
			buildOrigins: func() map[geoservices.LatLng]string {
				return map[geoservices.LatLng]string{
					geoservices.LatLng{40.736791925763455, -73.95519101851923}: "Saint Vitus Bar - Greenpoint",
					geoservices.LatLng{40.73622634374919, -73.95551867494544}:  "Le Fanfare - Greenpoint",
					geoservices.LatLng{40.72546946142794, -73.95175080861969}:  "Nassau Ave Station",
					geoservices.LatLng{40.726094743641546, -73.95230565534146}: "Peter Pan Donut & Pastry Shop",
					geoservices.LatLng{40.7217615358787, -73.95457096241172}:   "McCarren Park Tennis Courts",
				}
			},
			// Brooklyn City Dental
			destination: geoservices.LatLng{40.7263248173875, -73.95246643844668},
			expect: func(t *testing.T, out []geoservices.LatLng, originsMap map[geoservices.LatLng]string) {
				require.NotEmpty(t, out)

				// Convert output to human-readable places
				var ids []string
				for _, e := range out {
					id, ok := originsMap[e]
					if ok {
						ids = append(ids, id)
					}
				}

				require.Equal(t, []string{
					"Nassau Ave Station",
					"Peter Pan Donut & Pastry Shop",
					"McCarren Park Tennis Courts",
				}, ids[:3])

				// Sometimes ranking can be non-deterministic if, for example,
				// two points are equidistant from the target.
				require.Subset(t, ids[3:], []string{
					"Saint Vitus Bar - Greenpoint",
					"Le Fanfare - Greenpoint",
				})
			},
		},
	}
	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			ctx := context.Background()

			// Get inputs
			originsMap := tc.buildOrigins()
			var origins []geoservices.LatLng
			for k := range originsMap {
				origins = append(origins, k)
			}

			// Perform ranking
			out := Rank(ctx, origins, tc.destination)

			// Assertions
			tc.expect(t, out, originsMap)
		})
	}
}
