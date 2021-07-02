package busroutesapi

import (
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
)

// Bus describes http model of bus for api v1.
type Bus struct {
	ID     int64  `json:"id,omitempty"`
	CityID int    `json:"city_id,omitempty"`
	Num    string `json:"num"`

	City string `json:"city,omitempty"`
}

// Route describes http model of route for api v1.
type Route struct {
	BusID  int64 `json:"bus_id"`
	StopID int64 `json:"stop_id"`
	Step   int8  `json:"step"`
}

// RoutePoint describes a unit of route for a bus.
type RoutePoint struct {
	Step    int8   `json:"step"`
	Address string `json:"address"`
}

// RouteDetailed describes http model of detailed route for api v1.
type RouteDetailed struct {
	City   string       `json:"city"`
	Bus    string       `json:"bus"`
	Points []RoutePoint `json:"points"`
}

// RangeBusesResponse describes response for range of buses.
type RangeBusesResponse struct {
	Buses []*Bus `json:"items"`
	Total int64  `json:"total"`
}

// RangeBusesResponse describes response for buses.
type BusesResponse struct {
	Data  *RangeBusesResponse `json:"data,omitempty"`
	Error *ierr.APIError      `json:"error,omitempty"`
}

// RangeBusesResponse describes response for range of routes.
type RangeRoutesResponse struct {
	Routes []*RouteDetailed `json:"items"`
	Total  int64            `json:"total"`
}

// RangeBusesResponse describes response for routes.
type RoutesResponse struct {
	Data  *RangeRoutesResponse `json:"data,omitempty"`
	Error *ierr.APIError       `json:"error,omitempty"`
}
