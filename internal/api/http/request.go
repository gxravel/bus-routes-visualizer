package api

import (
	"fmt"
	"net/http"
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

	url := fmt.Sprintf("%s/buses?cities=%s&nums=%s", api, city, bus)
	return url, nil
}
