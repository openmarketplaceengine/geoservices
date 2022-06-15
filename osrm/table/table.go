package table

import (
	"encoding/json"
	"fmt"
	osrm "github.com/openmarketplaceengine/geoservices"
	"io/ioutil"
	"net/http"
	"net/url"
)

type Annotations string

const (
	Duration         Annotations = "duration"
	Distance         Annotations = "distance"
	DurationDistance Annotations = "duration,distance"
)

type Request struct {
	Coordinates  []osrm.Coordinate
	Origins      []int
	Destinations []int
	Annotations  Annotations
	//FallbackSpeed      *float64
	//FallbackCoordinate *osrm.Coordinate
	//ScaleFactor        *float64
}

type Response struct {
	Code         string          `json:"code"`
	Distances    [][]float64     `json:"distances"`
	Durations    [][]float64     `json:"durations"`
	Destinations []osrm.Waypoint `json:"destinations"`
	Sources      []osrm.Waypoint `json:"sources"`
}

// Table will return the fastest route between TableRequest Origins and each Destination.
// i.e. For worker to get ranked routes to available jobs.
// http request goes to http://project-osrm.org/docs/v5.23.0/api/#table-service endpoint.
func Table(c *http.Client, request Request) (*Response, error) {

	coords := "polyline(" + url.PathEscape(osrm.ToPolyline(request.Coordinates)) + ")"

	s := osrm.EncodeUrlParam("sources", request.Origins)
	d := osrm.EncodeUrlParam("destinations", request.Destinations)
	a := fmt.Sprintf("annotations=%s", request.Annotations)
	uri := fmt.Sprintf("https://router.project-osrm.org/table/v1/driving/%s?%s&%s&%s", coords, s, d, a)

	res, err := c.Get(uri)
	defer res.Body.Close()
	bytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read body: %w", err)
	}

	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusBadRequest {
		return nil, fmt.Errorf("unexpected response: %q status: %d", bytes, res.StatusCode)
	}

	var response Response
	err = json.Unmarshal(bytes, &response)
	if err != nil {
		return nil, fmt.Errorf("unmarshal response: %q error: %w", bytes, err)
	}

	return &response, nil
}
