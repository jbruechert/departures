package api

import (
	"context"
	"time"
)

// Remark is a hint or warning attached to a departure (e.g. bicycle conveyance).
type Remark struct {
	Type string `json:"type"`
	Code string `json:"code"`
	Text string `json:"text"`
}

// Operator is the company running a line.
type Operator struct {
	Type string `json:"type"`
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Line describes the transit line of a departure.
type Line struct {
	Type     string   `json:"type"`
	ID       string   `json:"id"`
	FahrtNr  string   `json:"fahrtNr"`
	Name     string   `json:"name"`
	Public   bool     `json:"public"`
	Mode     string   `json:"mode"`
	Product  string   `json:"product"`
	Operator Operator `json:"operator"`
}

// Departure is a single upcoming departure at a stop.
type Departure struct {
	TripID    string    `json:"tripId"`
	Stop      Location  `json:"stop"`
	When      time.Time `json:"when"`
	Direction string    `json:"direction"`
	Line      Line      `json:"line"`
	Remarks   []Remark  `json:"remarks"`
	Delay     int       `json:"delay"`
	Platform  string    `json:"platform"`
}

type departuresQuery struct {
	c        *Client
	id       string
	results  int
	language string
	duration int
}
type motisStopTimesResponse struct {
	StopTimes []motisStopTime `json:"stopTimes"`
	Place     motisPlace      `json:"place"`
}
type motisStopTime struct {
	Place motisPlace `json:"place"`

	Mode     string `json:"mode"`
	RealTime bool   `json:"realTime"`

	Headsign string `json:"headsign"`

	TripID string `json:"tripId"`

	RouteShortName string `json:"routeShortName"`
	RouteLongName  string `json:"routeLongName"`

	TripShortName string `json:"tripShortName"`
	DisplayName   string `json:"displayName"`

	AgencyID   string `json:"agencyId"`
	AgencyName string `json:"agencyName"`

	Track          string `json:"track"`
	ScheduledTrack string `json:"scheduledTrack"`

	Cancelled     bool `json:"cancelled"`
	TripCancelled bool `json:"tripCancelled"`
}

type motisPlace struct {
	Name string `json:"name"`
	ID   string `json:"stopId"`

	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`

	Arrival            time.Time `json:"arrival"`
	Departure          time.Time `json:"departure"`
	ScheduledArrival   time.Time `json:"scheduledArrival"`
	ScheduledDeparture time.Time `json:"scheduledDeparture"`
}

func (c *Client) Departures(id string) *departuresQuery {
	return &departuresQuery{
		c:        c,
		id:       id,
		results:  10,
		language: "en",
		duration: 60,
	}
}
func translateStopTimes(in motisStopTimesResponse) []Departure {
	out := make([]Departure, 0, len(in.StopTimes))

	for _, st := range in.StopTimes {
		when := st.Place.Departure
		if when.IsZero() {
			when = st.Place.ScheduledDeparture
		}
		out = append(out, Departure{
			TripID: st.TripID,

			Stop: Location{
				Type: "stop",
				ID:   st.Place.ID,
				Name: st.Place.Name,
				Location: struct {
					Type      string
					ID        string
					Latitude  float64
					Longitude float64
				}{
					Type:      "location",
					ID:        st.Place.ID,
					Latitude:  st.Place.Lat,
					Longitude: st.Place.Lon,
				},
			},

			When: when,

			Direction: st.Headsign,

			Line: Line{
				Type:    "line",
				ID:      st.RouteShortName,
				Name:    st.DisplayName,
				FahrtNr: st.TripShortName,
				Mode:    st.Mode,
				Product: st.Mode,
				Public:  true,
				Operator: Operator{
					Type: "operator",
					ID:   st.AgencyID,
					Name: st.AgencyName,
				},
			},
		})
	}
	return out
}

func (q *departuresQuery) Do(ctx context.Context) ([]Departure, error) {
	const u = "/api/v6/stoptimes?stopId=%s&n=%d&radius=200"

	var resp motisStopTimesResponse

	err := q.c.getJSON(
		ctx,
		&resp,
		u,
		q.id,
		q.results,
	)

	if err != nil {
		return nil, err
	}

	return translateStopTimes(resp), nil
}

func (q *departuresQuery) Results(r int) *departuresQuery {
	q.results = r
	return q
}

func (q *departuresQuery) Language(l string) *departuresQuery {
	q.language = l
	return q
}

func (q *departuresQuery) Duration(d int) *departuresQuery {
	q.duration = d
	return q
}
