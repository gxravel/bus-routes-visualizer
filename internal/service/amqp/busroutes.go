package amqp

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/api/http/handler/v1"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/service"

	"github.com/gxravel/bus-routes/pkg/rmq"
)

// BusRoutesService implements busroutes service interface.
type BusRoutesService struct {
	client *amqpClient
}

// NewBusRoutesService creates new busroutes client.
func NewBusRoutesService(client *rmq.Client) service.BusRoutes {
	return &BusRoutesService{
		client: newCustomClient(client),
	}
}

// GetRoutesDetailed makes request to the client and returns detailed routes:
func (s *BusRoutesService) GetRoutesDetailed(ctx context.Context, bus *httpv1.Bus) ([]*httpv1.RouteDetailed, error) {
	logger := log.FromContext(ctx).
		WithFields(
			"module", "GetRoutesDetailed",
			"bus", bus,
		)
	ctx = log.CtxWithLogger(ctx, logger)

	meta := &rmq.Meta{
		QName: "abcde",
	}

	rangeResponse := &httpv1.RangeRoutesResponse{}

	if err := s.client.processRequest(ctx, meta, bus, rangeResponse); err != nil {
		return nil, err
	}
	if rangeResponse == nil {
		return nil, nil
	}

	return rangeResponse.Routes, nil
}
