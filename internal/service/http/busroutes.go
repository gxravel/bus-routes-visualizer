package http

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/api/http/handler/v1"
	"github.com/gxravel/bus-routes-visualizer/internal/config"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/service"
)

// BusroutesService implements busroutes service interface.
// It uses http.
type BusroutesService struct {
	client *httpClient
	api    string
}

// NewBusroutesService creates new busroutes client.
func NewBusroutesService(logger log.Logger, conf *config.Config) service.Busroutes {
	var customClient = newCustomClient(
		withTimeout(conf),
		withUseInsecureTLS(conf),
	)

	return &BusroutesService{
		client: customClient,
		api:    conf.RemoteServices.BusroutesAPI,
	}
}

const (
	routeForBuses  = "/buses"
	routeForRoutes = "/routes/detailed"
)

// GetRoutesDetailed makes 2 requests to the API:
// 1) /buses for receiving buses ids
// 2) /routes/detailed for receiving routes.
func (s *BusroutesService) GetRoutesDetailed(ctx context.Context, bus *httpv1.Bus) ([]*httpv1.RouteDetailed, error) {
	url := fmt.Sprintf("%s?cities=%s&nums=%s", s.api+routeForBuses, bus.City, bus.Num)

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
