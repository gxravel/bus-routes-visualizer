package busroutesapi

import (
	"encoding/json"
	"net/http"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
	"github.com/gxravel/bus-routes-visualizer/internal/busroutesapi/v1"
	"github.com/gxravel/bus-routes-visualizer/internal/errors"
)

type ItemsType uint8

const (
	TypeBuses ItemsType = iota
	TypeRoutes
)

// GetItems makes request to the busroutes api
// and converts a RangeItemsResponse to the specified type.
func GetItems(url string, itemsType ItemsType) (int, interface{}, error) {
	response, err := http.Get(url)
	if err != nil {
		return 0, nil, err
	}

	var itemsResponse = &api.Response{}

	switch itemsType {
	case TypeBuses:
		itemsResponse.Data = &busroutesapi.RangeBusesResponse{}
	case TypeRoutes:
		itemsResponse.Data = &busroutesapi.RangeRoutesResponse{}
	default:
		itemsResponse.Data = &api.RangeItemsResponse{}
	}

	if err := processResponse(response, itemsResponse); err != nil {
		return 0, nil, err
	}

	if itemsResponse.Error != nil {
		itemsResponse.Data = nil
		return response.StatusCode, itemsResponse, errors.ErrProviderAPI
	}

	return 0, itemsResponse.Data, nil
}

func processResponse(r *http.Response, data interface{}) error {
	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return err
	}
	return nil
}
