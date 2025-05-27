package nats

import (
	"context"
	"fmt"

	"github.com/nats-io/nats.go"
	"github.com/ronnyp07/SportStream/worker/internal/metrics"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/infaestructure/log"
	"github.com/sts-solutions/base-code/cccorrelation"
	"github.com/sts-solutions/base-code/cclogger"
)

const (
	EventCorrelationIdKey = "correlation-id"
)

func Connect(ctx context.Context, host, port string, options ...nats.Options) (*nats.Conn, error) {
	url := natsConnectionString(host, port)

	return nats.Connect(url,
		nats.RetryOnFailedConnect(true),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			if err != nil {
				metrics.ReportNatsError(err)
				log.Logger().Error(ctx, fmt.Sprintf("nats client disconnected %v", err))
			}
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Logger().Info(ctx, "nats client reconnected")
		}),
		nats.ClosedHandler(func(_ *nats.Conn) {
			log.Logger().Info(ctx, "nats client closed")
		}),
		nats.ErrorHandler(func(_ *nats.Conn, _ *nats.Subscription, err error) {
			log.Logger().Error(ctx, "nats client connected error: %v", cclogger.LogField{
				Key:   "reason",
				Value: err.Error(),
			})
		}),
	)
}

func SetMsgCorrelationId(ctx context.Context, msg *nats.Msg) {
	if ok, id := cccorrelation.GetCorrelationId(ctx); ok {
		if msg.Header == nil {
			msg.Header = make(nats.Header)
		}

		msg.Header.Add(EventCorrelationIdKey, id)
	}
}

func natsConnectionString(host, port string) string {
	connectionString := fmt.Sprintf("nats://%s:%s", host, port)
	return connectionString
}
