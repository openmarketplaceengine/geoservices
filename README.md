# geoservices

[![Lines Of Code](https://tokei.rs/b1/github/openmarketplaceengine/geoservices?category=code)](https://github.com/openmarketplaceengine/geoservices)
[![Go Reference](https://pkg.go.dev/badge/github.com/openmarketplaceengine/geoservices.svg)](https://pkg.go.dev/github.com/openmarketplaceengine/geoservices)
[![Go Report Card](https://goreportcard.com/badge/github.com/openmarketplaceengine/geoservices)](https://goreportcard.com/report/github.com/openmarketplaceengine/geoservices)

This is a Go library that abstracts over various 3rd-party geoservice vendors:
* Google Maps
* OSMR
* GraphHopper (future)
* MapBox (future)

## Table of Contents
1. [Why use this?](#why-use-this\?)
2. [Getting started](#getting-started)
   1. [Ranking](#ranking)
   2. [Distance matrix](#distance-matrix)
   3. [Geocoding](#geocoding)

## Why use this?
🔌 **It's pluggable**. Use a 3rd-party vendor of your choice.

🧱 **It's modular**. Import only the functionality you need. ([Go 1.18 workspaces](https://go.dev/blog/get-familiar-with-workspaces) 
really help keep dependency trees minimal).

⚡ **It's performant**. We parallelize requests and partition distance 
matrices for you.

## Getting started

### Ranking
Install with:
```bash
go get github.com/openmarketplaceengine/geoservices/ranking
```

The [`ranking`](./ranking) package provides a function for getting the nearest 
points to a specified location.
```go
package main

import (
	"context"
	"github.com/openmarketplaceengine/geoservices"
	"github.com/openmarketplaceengine/geoservices/ranking"
)

func main() {
	origins := []geoservices.LatLng{
		{40.736791925763455, -73.95519101851923},
		{40.73622634374919, -73.95551867494544},
	}

	target := geoservices.LatLng{40.7263248173875, -73.95246643844668}

	var out []geoservices.LatLng
	out = ranking.Rank(context.Background(), origins, target)
}
```

This is a decent, stateless substitute for getting a rough distance-based 
ordering of points.

Output is non-deterministic, meaning the results may not necessarily have the 
same order given the same set of inputs. This might arise when two origin points
are equidistant from the target.

Performance is roughly 1ms per 100 locations on a 10-core Apple M1 Pro chip.

### Distance Matrix
The [`distance`](./distance) package abstracts over various Distance Matrix 
solutions.

For this example, we'll use Google Maps:
```bash
go get github.com/openmarketplaceengine/geoservices/distance/google
```

Then in Go:
```go
package main

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/geoservices"
	"github.com/openmarketplaceengine/geoservices/distance"
	"github.com/openmarketplaceengine/geoservices/distance/google"
	"googlemaps.github.io/maps"
	"os"
)

func main() {
	ctx := context.Background()

	origins := []geoservices.LatLng{
		// Prospect Park Tennis Center
		{40.651342660933835, -73.97032057647756},
		// Brooklyn Botanic Garden
		{40.66955251860222, -73.96234417342181},
	}

	// Barclays Center (Brooklyn)
	destination := geoservices.LatLng{40.68266462714047, -73.9754129316442}

	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		panic(err)
	}

	out, err := google.BetweenPoints(ctx, client, distance.BetweenPointsInput{
		Origins: origins,
		Destinations: []geoservices.LatLng{
			destination,
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
	fmt.Printf("Getting from %q to %q takes %v\n",
		"Prospect Park Tennis Center", // the 1st origin, hence Rows[0]
		"Barclays Center (Brooklyn)",  // the 1st destination, hence Elements[0]
		res.Rows[0].Elements[0].Duration,
	)
	fmt.Printf("Getting from %q to %q takes %v\n",
		"Brooklyn Botanic Garden",     // the 2nd origin, hence Rows[1]
		"Barclays Center (Brooklyn)",  // the 1st destination, hence Elements[0]
		res.Rows[1].Elements[0].Duration,
	)
}
```

In this example, we're feeding in 2 origins (tennis center + botanic garden) and
1 destination (Barclays Center).

The response will contain:
* a list of geocoded human-readable addresses for the origins (order matches input);
* a list of geocoded human-readable addresses for the destination (order matches input);
* a list of rows (per origin), where each row contains the pairwise distance/duration from said origin to each destination.

**The order of outputs will correspond to order of inputs**.

### Geocoding
The [`geocode`](./geocode) package abstracts over various geocoding solutions.

For this example, we'll use Google Maps:
```bash
go get github.com/openmarketplaceengine/geoservices/geocode/google
```

Then in Go:
```go
package main

import (
	"context"
	"fmt"
	"github.com/openmarketplaceengine/geoservices"
	"github.com/openmarketplaceengine/geoservices/geocode/google"
	"googlemaps.github.io/maps"
	"os"
)

func main() {
	ctx := context.Background()
	target := geoservices.LatLng{40.7263248173875, -73.95246643844668}
	client, err := maps.NewClient(maps.WithAPIKey(os.Getenv("GOOGLE_MAPS_API_KEY")))
	if err != nil {
		panic(err)
	}
	geocoder := google.NewGeocoder(client)
	res, err := geocoder.ReverseGeocode(ctx, target)
	if err != nil {
		panic(err)
	}
	fmt.Println(res)
}
```