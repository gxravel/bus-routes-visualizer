package amqp

import (
	"context"

	httpv1 "github.com/gxravel/bus-routes-visualizer/internal/api/http/handler/v1"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"
	"github.com/gxravel/bus-routes-visualizer/internal/service"

	"github.com/gxravel/bus-routes/pkg/rmq"
	"github.com/streadway/amqp"
)

// BusroutesService implements busroutes service interface.
// It uses amqp.
type BusroutesService struct {
	*amqpClient
	consumer *rmq.Consumer
}

// NewBusroutesService creates new busroutes client.
func NewBusroutesService(ctx context.Context, publisher *rmq.Publisher, consumer *rmq.Consumer) (service.Busroutes, error) {
	s := &BusroutesService{
		newCustomClient(publisher),
		consumer,
	}

	if err := s.handleRPCReplies(context.TODO()); err != nil {
		return nil, err
	}

	return s, nil
}

// GetRoutesDetailed makes request to the client and returns detailed routes:
func (s *BusroutesService) GetRoutesDetailed(ctx context.Context, bus *httpv1.Bus) ([]*httpv1.RouteDetailed, error) {
	logger := log.FromContext(ctx).
		WithModule("GetRoutesDetailed").
		WithField("bus", bus)
	ctx = log.CtxWithLogger(ctx, logger)

	rangeResponse := &httpv1.RangeRoutesResponse{}

	qname, free := s.consumer.GetFreeQueue()
	defer free()

	meta := rmq.GetMetaDetailedRoutesRPC()
	meta.QName = qname

	if err := s.processRequest(ctx, meta, bus, rangeResponse); err != nil {
		return nil, err
	}
	if rangeResponse == nil {
		return nil, nil
	}

	return rangeResponse.Routes, nil
}

func (s *BusroutesService) handleRPCReplies(ctx context.Context) error {
	for _, qname := range s.consumer.ListAllQueues() {
		delivery, err := s.consumer.Consume(qname, true, true)
		if err != nil {
			return err
		}

		go func(delivery <-chan amqp.Delivery) {
			for {
				select {
				case <-ctx.Done():
					return

				case message := <-delivery:
					s.deliveries[message.CorrelationId] <- &message
				}
			}
		}(delivery)
	}

	return nil
}
