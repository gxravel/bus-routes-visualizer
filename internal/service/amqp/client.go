package amqp

import (
	"context"
	"errors"
	"time"

	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"

	"github.com/google/uuid"
	"github.com/gxravel/bus-routes/pkg/rmq"
	amqpv1 "github.com/gxravel/bus-routes/pkg/rmq/v1"
	"github.com/streadway/amqp"
)

const (
	defaultTimeout = time.Second * 5
)

// amqpClient wraps rmq.Publisher to publish requests.
type amqpClient struct {
	publisher *rmq.Publisher

	// deliveries is to bring delivery to the particular correlation ID.
	deliveries map[string]chan *amqp.Delivery
}

// newCustomClient creates new instance of amqpClient
func newCustomClient(
	publisher *rmq.Publisher,
) *amqpClient {

	return &amqpClient{
		publisher: publisher,

		deliveries: make(map[string]chan *amqp.Delivery),
	}
}

// processRequest processes a request by calling RPC, processing response and logging the data.
func (c *amqpClient) processRequest(ctx context.Context, meta *rmq.Meta, body, result interface{}) error {
	publisher, free := c.publisher.WithFreeChannel()
	defer free()

	client := newCustomClient(publisher)

	logger := log.FromContext(ctx)

	defer func(start time.Time) {
		logger.
			WithFields(
				"duration", time.Since(start),
				"meta", meta,
			).
			Debug("processed amqp request")
	}(time.Now())

	meta.CorrID = uuid.New().String()

	c.deliveries[meta.CorrID] = make(chan *amqp.Delivery)
	defer delete(c.deliveries, meta.CorrID)

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	response := &amqpv1.Response{Data: result}
	if err := client.callRPC(ctx, meta, c.deliveries[meta.CorrID], body, response); err != nil {
		logger = logger.WithErr(err)
		return err
	}

	logger = logger.WithField("response", response)

	if err := client.processResponse(response); err != nil {
		logger = logger.WithErr(err)
		return err
	}

	return nil
}

// callRPC calls RPC with message body.
// Waits an answer and writes it to the response.
func (c *amqpClient) callRPC(ctx context.Context, meta *rmq.Meta, delivery chan *amqp.Delivery, body, response interface{}) error {
	messageBody, err := rmq.ConvertToMessage(body)
	if err != nil {
		return err
	}

	if err := c.publisher.Produce(meta, messageBody); err != nil {
		return err
	}

	select {
	case message := <-delivery:
		if err := rmq.TranslateMessage(message.Body, response); err != nil {
			return err
		}

	case <-ctx.Done():
		return errors.New("context done")
	}

	return nil
}

// processResponse processses a response by making the checks and handling response error.
func (c *amqpClient) processResponse(response *amqpv1.Response) error {
	if response == nil {
		return nil
	}

	if response.Error != nil {
		return ierr.NewProviderAPIError(
			response.Error.Reason.Err+": "+response.Error.Reason.Message,
			response.Error.Code,
		)
	}

	return nil
}
