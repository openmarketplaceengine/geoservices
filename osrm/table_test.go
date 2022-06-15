package osrm

import (
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestTable(t *testing.T) {
	request := TableRequest{
		Coordinates:        []Coordinate{{40.751753, -73.980052}, {40.794156, -73.962662}},
		Origins:            []int{0},
		Destinations:       []int{0, 1},
		Annotations:        DurationDistance,
		FallbackSpeed:      nil,
		FallbackCoordinate: nil,
		ScaleFactor:        nil,
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

	//bytes, err := json.Marshal(response)
	//require.NoError(t, err)
	//fmt.Printf("%s", bytes)

	//coords2 := [][]float64{}
	//for _, route := range resp.Routes {
	//	for _, leg := range route.Legs {
	//		for _, step := range leg.Steps {
	//			g := step.Geometry
	//
	//			pointSet := g.PointSet
	//			for _, point := range pointSet {
	//				coords2 = append(coords2, []float64{point.Lat(), point.Lng()})
	//			}
	//		}
	//	}
	//}
	//fmt.Println(string(polyline.EncodeCoords(coords2)))
	//coords := [][]float64{
	//	{38.5, -120.2},
	//	{40.7, -120.95},
	//	{43.252, -126.453},
	//}
	//fmt.Println(string(polyline.EncodeCoords(coords)))
}
