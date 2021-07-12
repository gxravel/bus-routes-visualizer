package amqp

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/api/http/handler/v1"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/service"

	"github.com/gxravel/bus-routes/pkg/rmq"
)

// BusroutesService implements busroutes service interface.
// It uses amqp.
type BusroutesService struct {
	client *amqpClient
}

// NewBusroutesService creates new busroutes client.
func NewBusroutesService(client *rmq.Client) service.Busroutes {
	return &BusroutesService{
		client: newCustomClient(client),
	}
}

// GetRoutesDetailed makes request to the client and returns detailed routes:
func (s *BusroutesService) GetRoutesDetailed(ctx context.Context, bus *httpv1.Bus) ([]*httpv1.RouteDetailed, error) {
	logger := log.FromContext(ctx).
		WithFields(
			"module", "GetRoutesDetailed",
			"bus", bus,
		)
	ctx = log.CtxWithLogger(ctx, logger)

	rangeResponse := &httpv1.RangeRoutesResponse{}

	if err := s.client.processRequest(ctx, rmq.MetaDetailedRoutesRPC, bus, rangeResponse); err != nil {
		return nil, err
	}
	if rangeResponse == nil {
		return nil, nil
	}

	return rangeResponse.Routes, nil
}
