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
	*amqpClient
}

// NewBusroutesService creates new busroutes client.
func NewBusroutesService(publisher *rmq.Publisher, consumer *rmq.Consumer) service.Busroutes {
	return &BusroutesService{
		amqpClient: newCustomClient(publisher, consumer),
	}
}

// GetRoutesDetailed makes request to the client and returns detailed routes:
func (s *BusroutesService) GetRoutesDetailed(ctx context.Context, bus *httpv1.Bus) ([]*httpv1.RouteDetailed, error) {
	logger := log.FromContext(ctx).
		WithModule("GetRoutesDetailed").
		WithField("bus", bus)
	ctx = log.CtxWithLogger(ctx, logger)

	rangeResponse := &httpv1.RangeRoutesResponse{}

	if err := s.processRequest(ctx, rmq.GetMetaDetailedRoutesRPC(), bus, rangeResponse); err != nil {
		return nil, err
	}
	if rangeResponse == nil {
		return nil, nil
	}

	return rangeResponse.Routes, nil
}
