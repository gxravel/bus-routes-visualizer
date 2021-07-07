package v1

import (
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
)

// Response describes http range itmes response for api v1.
type RangeItemsResponse struct {
	Items interface{} `json:"items"`
	Total int64       `json:"total"`
}

// Response describes http response for api v1.
type Response struct {
	Data  interface{}    `json:"data,omitempty"`
	Error *ierr.APIError `json:"error,omitempty"`
}

// Response describes http model of permission for api v1.
type Permission struct {
	UserID  int64       `json:"user_id"`
	Actions interface{} `json:"actions"`
}
