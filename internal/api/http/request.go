package api

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gxravel/bus-routes-visualizer/internal/dataprovider"
	"github.com/pkg/errors"
)

// ParseQueryParam parses query params for specific field.
func ParseQueryParam(r *http.Request, field string) (string, error) {
	q := r.URL.Query()

	param := q.Get(field)
	if param == "" {
		return "", nil
	}

	return param, nil
}

// ParseQueryParams parses query params for specific field.
func ParseQueryParams(r *http.Request, field string) ([]string, error) {
	q := r.URL.Query()
	params := q[field]

	if len(params) == 0 {
		return nil, nil
	}

	return params, nil
}

// ParseQueryInt64Slice parses query []int64 for specific field.
func ParseQueryInt64Slice(r *http.Request, field string) ([]int64, error) {
	q := r.URL.Query()
	params := q[field]

	if len(params) == 0 {
		return nil, nil
	}

	var vals []int64

	for _, p := range params {
		slice := strings.Split(p, ",")
		if vals == nil {
			vals = make([]int64, 0, len(slice))
		}

		for _, s := range slice {
			if s == "" {
				continue
			}
			val, err := strconv.ParseInt(s, 10, 64)
			if err != nil {
				return nil, errors.Errorf("can't parse %v to int", s)
			}
			vals = append(vals, val)
		}
	}
	return vals, nil
}

// ParseGraphsRequest returns the url for further request to get bus ids.
func ParseGraphsRequest(r *http.Request, api string) (string, error) {
	bus, err := ParseQueryParam(r, "bus")
	if err != nil {
		return "", err
	}

	city, err := ParseQueryParam(r, "city")
	if err != nil {
		return "", err
	}

	url := fmt.Sprintf("%s?cities=%s&nums=%s", api, city, bus)
	return url, nil
}

// ParsePermissionsFilter parses user_ids and actions, and returns the filter.
func ParsePermissionsFilter(r *http.Request) (*dataprovider.PermissionFilter, error) {
	ids, err := ParseQueryInt64Slice(r, "user_ids")
	if err != nil {
		return nil, err
	}

	actions, err := ParseQueryParams(r, "actions")
	if err != nil {
		return nil, err
	}

	return dataprovider.
		NewPermissionFilter().
		ByUserIDs(ids...).
		ByActions(actions...), nil
}
