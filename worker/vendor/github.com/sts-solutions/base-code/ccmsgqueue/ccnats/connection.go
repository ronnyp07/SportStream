package ccnats

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"github.com/sts-solutions/base-code/ccmsgqueue"
)

const (
	ConnectionError         = "ConnectionError"
	JetStreamError          = "JetStreamError"
	ClientDisconnectedError = "ClientDisconnectedError"
	GenericError            = "GenericError"
)

type connection struct {
	port     int
	host     string
	tracer   Tracer
	logger   ccmsgqueue.Logger
	conn     *nats.Conn
	js       jetstream.JetStream
	metricts ccmsgqueue.ConnectionMetrics
}

func (c *connection) Connect() (err error) {
	c.conn, err = c.connect()
	if err != nil {
		c.metricts.ErrorInc(ConnectionError)
		return errors.Wrap(err, "connecting to nats server")
	}

	c.js, err = jetstream.New(c.conn)
	if err != nil {
		c.metricts.ErrorInc(JetStreamError)
		return errors.Wrap(err, "getting JetStream")
	}
	c.tracer.setConnectionTracer(c.conn)
	return nil
}

func (c *connection) IsConnected() bool {
	return c.conn != nil && c.conn.IsConnected()
}

func (c *connection) Close() {
	c.conn.Close()
}

func (c *connection) connect() (*nats.Conn, error) {
	url := fmt.Sprintf("nats://%s:%d", c.host, c.port)
	ctx := context.Background()

	opts := []nats.Option{
		nats.RetryOnFailedConnect(true),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				c.logger.LogError(ctx, fmt.Sprintf("nats client disconnected. Error: %v", err.Error()))
				c.metricts.ErrorInc(ClientDisconnectedError)
			}
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			c.logger.LogInfo(ctx, "nats client closed")
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			c.logger.LogError(ctx, fmt.Sprintf("nats Error: %v", err.Error()))
			c.metricts.ErrorInc(GenericError)
		}),
	}

	c.logger.LogInfo(ctx, fmt.Sprintf("connecting to nats. URL: %s", url))
	return nats.Connect(url, opts...)
}
