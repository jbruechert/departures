package api

import (
	"context"
	"fmt"
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
	duration int
	results  int
	language string
}

// Departures builds a query for the upcoming departures at the stop with the given ID.
func (c *Client) Departures(id string) *departuresQuery {
	return &departuresQuery{
		c:        c,
		id:       id,
		duration: 60,
		results:  0, // 0 means "no limit" — let duration decide.
		language: "en",
	}
}

// departuresResponse wraps the v6 departures payload, which is an object
// ({ "departures": [...], "realtimeDataUpdatedAt": ... }) rather than a bare array.
type departuresResponse struct {
	Departures []Departure `json:"departures"`
}

func (q *departuresQuery) Do(ctx context.Context) ([]Departure, error) {
	const u = "/stops/%s/departures?duration=%d&language=%s&pretty=false%s"

	// Only cap the number of results when asked; otherwise let duration decide,
	// matching the API default (an unset results means "no limit").
	results := ""
	if q.results > 0 {
		results = fmt.Sprintf("&results=%d", q.results)
	}

	var resp departuresResponse
	err := q.c.getJSON(ctx, &resp, u, q.id, q.duration, q.language, results)
	if err != nil {
		return nil, err
	}

	return resp.Departures, nil
}

func (q *departuresQuery) Duration(d int) *departuresQuery {
	q.duration = d
	return q
}

func (q *departuresQuery) Results(r int) *departuresQuery {
	q.results = r
	return q
}

func (q *departuresQuery) Language(l string) *departuresQuery {
	q.language = l
	return q
}
