package api

import (
	"context"
	"net/url"
)

type Location struct {
	Type     string
	ID       string
	Name     string
	Location struct {
		Type      string
		ID        string
		Latitude  float64
		Longitude float64
	}
	Products    map[string]bool
	StationDHID string
}

type locationsQuery struct {
	c            *Client
	query        string
	fuzzy        bool
	results      int
	stops        bool
	addresses    bool
	poi          bool
	linesOfStops bool
	language     string
}

func (c *Client) Locations(query string) *locationsQuery {
	return &locationsQuery{
		c:            c,
		query:        query,
		fuzzy:        true,
		results:      10,
		stops:        true,
		addresses:    true,
		poi:          true,
		linesOfStops: false,
		language:     "en",
	}
}

type motisLocation struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Type string `json:"type"`

	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`

	Place string `json:"place"`

	Stops []motisLocationStop `json:"stops"`
}

type motisLocationStop struct {
	ID string `json:"id"`

	Name string `json:"name"`

	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

func translateLocation(l motisLocation) Location {
	var loc Location

	loc.Type = l.Type
	loc.ID = l.ID
	loc.Name = l.Name

	loc.Location.Type = "coordinates"
	loc.Location.ID = l.ID
	loc.Location.Latitude = l.Lat
	loc.Location.Longitude = l.Lon

	return loc
}

func translateLocations(resp []motisLocation) []Location {
	locs := make([]Location, 0, len(resp))

	for _, l := range resp {
		locs = append(locs, translateLocation(l))
	}

	return locs
}

func (q *locationsQuery) Do(ctx context.Context) ([]Location, error) {
	const u = "/api/v1/geocode?text=%s&limit=%d"

	var resp []motisLocation

	err := q.c.getJSON(
		ctx,
		&resp,
		u,
		url.QueryEscape(q.query),
		q.results,
	)

	if err != nil {
		return nil, err
	}

	return translateLocations(resp), nil
}

func (q *locationsQuery) Fuzzy(f bool) *locationsQuery {
	q.fuzzy = f
	return q
}

func (q *locationsQuery) Results(r int) *locationsQuery {
	q.results = r
	return q
}

func (q *locationsQuery) Stops(s bool) *locationsQuery {
	q.stops = s
	return q
}

func (q *locationsQuery) Addresses(a bool) *locationsQuery {
	q.addresses = a
	return q
}

func (q *locationsQuery) POI(p bool) *locationsQuery {
	q.poi = p
	return q
}

func (q *locationsQuery) LinesOfStops(l bool) *locationsQuery {
	q.linesOfStops = l
	return q
}

func (q *locationsQuery) Language(l string) *locationsQuery {
	q.language = l
	return q
}
