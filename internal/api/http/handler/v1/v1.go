package v1

import (
	"github.com/gxravel/bus-routes-visualizer/internal/model"
)

// Response describes http range itmes response for api v1.
type RangeItemsResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

// Response describes http response for api v1.
type Response struct {
	Data  interface{} `json:"data,omitempty"`
	Error *APIError   `json:"error,omitempty"`
}

// APIReason describes http model of error reason.
type APIReason struct {
	RType   string `json:"type"`
	Err     string `json:"error"`
	Message string `json:"message,omitempty"`
}

// APIError describes http model of error.
type APIError struct {
	Reason *APIReason `json:"reason"`
}

// Response describes http model of permission for api v1.
type Permission struct {
	UserID  int64       `json:"user_id"`
	Actions interface{} `json:"actions"`
}

// User describes http model of user for api v1.
type User struct {
	ID   int64          `json:"id"`
	Type model.UserType `json:"type"`
}

// Bus describes http model of bus for api v1.
type Bus struct {
	ID     int64  `json:"id,omitempty"`
	CityID int    `json:"city_id,omitempty"`
	Num    string `json:"num"`

	City string `json:"city,omitempty"`
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
	Error *APIError           `json:"error,omitempty"`
}

// RangeBusesResponse describes response for range of routes.
type RangeRoutesResponse struct {
	Routes []*RouteDetailed `json:"items"`
	Total  int64            `json:"total"`
}

// RangeBusesResponse describes response for routes.
type RoutesResponse struct {
	Data  *RangeRoutesResponse `json:"data,omitempty"`
	Error *APIError            `json:"error,omitempty"`
}
