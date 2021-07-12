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
	*rmq.Client
}

// newCustomClient creates new instance of amqpClient
func newCustomClient(client *rmq.Client) *amqpClient {
	c := &amqpClient{
		Client: client,
	}

	return c
}

// CallRPC calls RPC with message body.
// Waits an answer and writes it to the response.
func (c *amqpClient) CallRPC(ctx context.Context, meta *rmq.Meta, body, response interface{}) error {
	messageBody, err := c.ConvertToMessage(body)
	if err != nil {
		return err
	}

	delivery, corrID, err := c.Client.CallRPC(meta.QName, messageBody)
	if err != nil {
		return errors.Wrapf(err, "can not call RPC with meta: %v", meta)
	}

	for {
		select {
		case <-ctx.Done():
			return nil

		case message := <-delivery:
			if message.CorrelationId != corrID {
				continue
			}

			if err := c.TranslateMessage(message.Body, response); err != nil {
				return err
			}

			return nil
		}
	}
}

// processRequest processes a request by calling RPC, processing response and logging the data.
func (c *amqpClient) processRequest(ctx context.Context, meta *rmq.Meta, body, result interface{}) error {
	logger := log.FromContext(ctx).WithField("meta", meta)

	defer func(start time.Time) {
		logger.WithField("duration", time.Since(start)).Debug("processed amqp request")
	}(time.Now())

	ctx, cancel := context.WithTimeout(ctx, defaultTimeout)
	defer cancel()

	response := &amqpv1.Response{Data: result}
	if err := c.CallRPC(ctx, meta, body, response); err != nil {
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
