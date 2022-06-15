package osrm

//Coordinate is Longitude[0] Latitude[1] pair
type Coordinate [2]float64

// Waypoint describes source or destination
type Waypoint struct {
	Hint     string     `json:"hint"`
	Distance float64    `json:"distance"`
	Name     string     `json:"name"`
	Location Coordinate `json:"location"`
}
