package busroutesapi

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes-visualizer/internal/busroutesapi/v1"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/visualizercontext"
)

const (
	RouteForBuses  = "/buses"
	RouteForRoutes = "/routes/detailed"
)

type ItemsType uint8

const (
	TypeBuses ItemsType = iota
	TypeRoutes
)

// GetRoutesDetailed makes 2 requests to the API:
// 1) /buses for receiving buses ids
// 2) /routes/detailed for receiving routes.
func GetRoutesDetailed(ctx context.Context, api string, url string) ([]*busroutesapi.RouteDetailed, error) {
	logger := log.FromContext(ctx).WithStr("url", url)
	logger.Debug("going for buses")

	setToken := false

	data, err := getItems(ctx, url, setToken, TypeBuses)
	if err != nil {
		return nil, err
	}
	buses, ok := data.(*busroutesapi.RangeBusesResponse)
	if !ok {
		logger.WithField("data", data).Error("expect RangeBusesResponse")
		return nil, ierr.ErrInternalServer
	}
	if buses.Total == 0 {
		return nil, nil
	}

	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(api + RouteForRoutes + "?bus_ids=")
	for _, bus := range buses.Buses {
		urlBuilder.WriteString(strconv.FormatInt(bus.ID, 10))
		urlBuilder.WriteString(",")
	}
	url = urlBuilder.String()[:urlBuilder.Len()-1]

	logger = log.FromContext(ctx).WithStr("url", url)
	logger.Debug("going for routes")

	setToken = true

	data, err = getItems(ctx, url, setToken, TypeRoutes)
	if err != nil {
		return nil, err
	}
	routes, ok := data.(*busroutesapi.RangeRoutesResponse)
	if !ok {
		logger.WithField("data", data).Error("expect RangeRoutesResponse")
		return nil, err
	}
	if routes.Total == 0 {
		return nil, nil
	}

	return routes.Routes, nil
}

// getItems makes request to the busroutes api
// and converts a RangeItemsResponse to the specified type.
func getItems(ctx context.Context, url string, setToken bool, itemsType ItemsType) (interface{}, error) {
	request, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if setToken {
		setAuthToken(ctx, request)
	}

	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, err
	}

	var itemsResponse = &httpv1.Response{}

	switch itemsType {
	case TypeBuses:
		itemsResponse.Data = &busroutesapi.RangeBusesResponse{}
	case TypeRoutes:
		itemsResponse.Data = &busroutesapi.RangeRoutesResponse{}
	default:
		itemsResponse.Data = &httpv1.RangeItemsResponse{}
	}

	if err := processResponse(response, itemsResponse); err != nil {
		return nil, err
	}

	if itemsResponse.Error != nil {
		return nil, ierr.NewTypedError(
			itemsResponse.Error.Reason.RType,
			ierr.NewProviderAPIError(
				itemsResponse.Error.Reason.Err,
				response.StatusCode,
			),
		)
	}

	return itemsResponse.Data, nil
}

func processResponse(r *http.Response, data interface{}) error {
	return json.NewDecoder(r.Body).Decode(data)
}

func setAuthToken(ctx context.Context, r *http.Request) {
	r.Header.Set("Authorization", "Bearer "+visualizercontext.GetToken(ctx))
}
