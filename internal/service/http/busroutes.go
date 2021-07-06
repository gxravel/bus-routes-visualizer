package http

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/gxravel/bus-routes-visualizer/internal/config"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/service"
	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/service/http/v1"
)

// BusRoutesService implements busroutes service interface.
type BusRoutesService struct {
	client *httpClient
	api    string
}

// NewBusRoutesService creates new busroutes client.
func NewBusRoutesService(logger log.Logger, conf *config.Config) service.BusRoutes {
	var customClient = newCustomClient(
		withTimeout(conf),
		withUseInsecureTLS(conf),
	)

	return &BusRoutesService{
		client: customClient,
		api:    conf.API.BusRoutes,
	}
}

const (
	RouteForBuses  = "/buses"
	routeForRoutes = "/routes/detailed"
)

// GetRoutesDetailed makes 2 requests to the API:
// 1) /buses for receiving buses ids
// 2) /routes/detailed for receiving routes.
func (s *BusRoutesService) GetRoutesDetailed(ctx context.Context, url string) ([]*httpv1.RouteDetailed, error) {
	logger := log.FromContext(ctx).WithStr("url", url)
	logger.Debug("going for buses")

	busesResp := &httpv1.BusesResponse{}
	if err := s.client.processRequest(ctx, http.MethodGet, url, nil, busesResp); err != nil {
		return nil, err
	}
	if busesResp.Data.Total == 0 {
		return nil, nil
	}
	buses := busesResp.Data.Buses

	urlBuilder := strings.Builder{}
	urlBuilder.WriteString(s.api + routeForRoutes + "?bus_ids=")
	for _, bus := range buses {
		urlBuilder.WriteString(strconv.FormatInt(bus.ID, 10))
		urlBuilder.WriteString(",")
	}
	url = urlBuilder.String()[:urlBuilder.Len()-1]

	logger = log.FromContext(ctx).WithStr("url", url)
	logger.Debug("going for routes")

	routesResp := &httpv1.RoutesResponse{}
	if err := s.client.processRequest(ctx, http.MethodGet, url, nil, routesResp); err != nil {
		return nil, err
	}
	if routesResp.Data.Total == 0 {
		return nil, nil
	}

	return routesResp.Data.Routes, nil
}
