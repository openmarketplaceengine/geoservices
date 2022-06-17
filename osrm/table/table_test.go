package table

import (
	"github.com/openmarketplaceengine/geoservices/osrm"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestTable(t *testing.T) {
	request := Request{
		Coordinates:  []osrm.Coordinate{{40.751753, -73.980052}, {40.794156, -73.962662}},
		Origins:      []int{0},
		Destinations: []int{0, 1},
		Annotations:  DurationDistance,
	}

	response, err := Table(&http.Client{}, request)
	require.NoError(t, err)
	require.Equal(t, "Ok", response.Code)

	require.Len(t, response.Distances, 1)
	require.Len(t, response.Distances[0], 2)
	require.Len(t, response.Durations, 1)
	require.Len(t, response.Durations[0], 2)

	require.Len(t, response.Destinations, 2)
	require.Len(t, response.Sources, 1)

	require.Greater(t, response.Distances[0][1], float64(0))
	require.Greater(t, response.Durations[0][1], float64(0))
}

func TestTableMultipleOrigins(t *testing.T) {
	request := Request{
		Coordinates: []osrm.Coordinate{
			{40.751753, -73.980052},
			{40.794156, -73.962662},
			{40.785965, -73.952794},
		},
		Origins:      []int{0, 1},
		Destinations: []int{0, 1, 2},
		Annotations:  DurationDistance,
	}

	response, err := Table(&http.Client{}, request)
	require.NoError(t, err)
	require.Equal(t, "Ok", response.Code)

	require.Len(t, response.Destinations, 3)
	require.Len(t, response.Sources, 2)

	require.Len(t, response.Distances, 2)
	require.Len(t, response.Distances[0], 3)
	require.Len(t, response.Distances[1], 3)
	require.Len(t, response.Durations, 2)
	require.Len(t, response.Durations[0], 3)
	require.Len(t, response.Durations[1], 3)

}
