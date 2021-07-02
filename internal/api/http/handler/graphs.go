package handler

import (
	"net/http"
	"strconv"
	"strings"

	api "github.com/gxravel/bus-routes-visualizer/internal/api/http"
	busroutesapi "github.com/gxravel/bus-routes-visualizer/internal/busroutesapi"
	busroutesapiV1 "github.com/gxravel/bus-routes-visualizer/internal/busroutesapi/v1"
	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
)

func (s *Server) getGraph(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	url, err := api.ParseGraphsRequest(r, s.busroutesAPI)
	if err != nil {
		api.RespondError(ctx, w, err)
		return
	}

	logger := s.logger.WithStr("url", url)
	logger.Debug("going to buses")

	code, data, err := busroutesapi.GetItems(url, busroutesapi.TypeBuses)
	if err != nil {
		if err == ierr.ErrProviderAPI {
			api.RespondJSON(ctx, w, code, data)
			return
		}
		api.RespondError(ctx, w, err)
		return
	}
	buses, ok := data.(*busroutesapiV1.RangeBusesResponse)
	if !ok {
		logger.WithField("data", data).Error("expect RangeBusesResponse")
		api.RespondError(ctx, w, ierr.ErrInternalServer)
		return
	}
	if buses.Total == 0 {
		api.RespondEmptyItems(ctx, w)
		return
	}

	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(s.busroutesAPI + "/routes/detailed?bus_ids=")
	for _, bus := range buses.Buses {
		urlBuilder.WriteString(strconv.FormatInt(bus.ID, 10))
		urlBuilder.WriteString(",")
	}
	url = urlBuilder.String()[:urlBuilder.Len()-1]

	logger = s.logger.WithStr("url", url)
	logger.Debug("going to routes")

	code, data, err = busroutesapi.GetItems(url, busroutesapi.TypeRoutes)
	if err != nil {
		if err == ierr.ErrProviderAPI {
			api.RespondJSON(ctx, w, code, data)
			return
		}
		api.RespondError(ctx, w, err)
		return
	}
	routes, ok := data.(*busroutesapiV1.RangeRoutesResponse)
	if !ok {
		logger.WithField("data", data).Error("expect RangeRoutesResponse")
		api.RespondError(ctx, w, err)
		return
	}
	if routes.Total == 0 {
		api.RespondEmptyItems(ctx, w)
		return
	}

	api.RespondDataOK(ctx, w, routes)
}
