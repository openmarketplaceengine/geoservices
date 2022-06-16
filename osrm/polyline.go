package osrm

import "github.com/twpayne/go-polyline"

//ToPolyline will return coordinates google polyline encoded
//see https://developers.google.com/maps/documentation/utilities/polylineutility
//Output example: _p~iF~ps|U_ulLnnqC_mqNvxq`@
func ToPolyline(coordinates []Coordinate) string {
	var coords = make([][]float64, len(coordinates))
	for i, coordinate := range coordinates {
		coords[i] = []float64{coordinate[0], coordinate[1]}
	}
	return string(polyline.EncodeCoords(coords))
}
