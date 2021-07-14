package amqp

import (
	"context"
	"time"

	ierr "github.com/gxravel/bus-routes-visualizer/internal/errors"
	log "github.com/gxravel/bus-routes-visualizer/internal/logger"

	"github.com/gxravel/bus-routes/pkg/rmq"
	amqpv1 "github.com/gxravel/bus-routes/pkg/rmq/v1"
	"github.com/pkg/errors"
)

const (
	defaultTimeout = time.Second * 5
)

// amqpClient wraps rmq.Client to interact with RabbitMQ.
type amqpClient struct {
	publisher *rmq.Publisher
	consumer  *rmq.Consumer
}

// newCustomClient creates new instance of amqpClient
func newCustomClient(
	publisher *rmq.Publisher,
	consumer *rmq.Consumer,
) *amqpClient {

	return &amqpClient{
		publisher: publisher,
		consumer:  consumer,
	}
}

// processRequest processes a request by calling RPC, processing response and logging the data.
func (c *amqpClient) processRequest(ctx context.Context, meta *rmq.Meta, body, result interface{}) error {
	c.publisher.UseFreeChannel()
	c.consumer.UseFreeChannel()

	defer func() {
		c.publisher.FreeChannel()
		c.consumer.FreeChannel()
	}()

	logger := log.FromContext(ctx)

	defer func(start time.Time) {
		logger.
			WithFields(
				"duration", time.Since(start),
				"meta", meta,
			).
			Debug("processed amqp request")
	}(time.Now())

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	response := &amqpv1.Response{Data: result}
	if err := c.callRPC(ctx, meta, body, response); err != nil {
		logger = logger.WithErr(err)
		return err
	}

	logger = logger.WithField("response", response)

	if err := c.processResponse(response); err != nil {
		logger = logger.WithErr(err)
		return err
	}

	return nil
}

// callRPC calls RPC with message body.
// Waits an answer and writes it to the response.
func (c *amqpClient) callRPC(ctx context.Context, meta *rmq.Meta, body, response interface{}) error {
	messageBody, err := rmq.ConvertToMessage(body)
	if err != nil {
		return err
	}

	if err := c.publisher.CallRPC(meta, messageBody); err != nil {
		return errors.Wrap(err, "failed to call rpc")
	}

	delivery, err := c.consumer.Consume(meta.QName)
	if err != nil {
		return errors.Wrap(err, "failed to consume")
	}

	for {
		select {
		case <-ctx.Done():
			return errors.New("context is done")

		case message := <-delivery:
			if message.CorrelationId != meta.CorrID {
				continue
			}

			if err := rmq.TranslateMessage(message.Body, response); err != nil {
				return err
			}

			return nil
		}
	}
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
