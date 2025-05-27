package app

import (
	"context"
	"fmt"
	"time"

	"emperror.dev/errors"
	"github.com/ronnyp07/SportStream/worker/internal/metrics"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/config"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/infaestructure/log"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/otel"
	ccmetricsnats "github.com/sts-solutions/base-code/ccmetrics/ccmsgqueue/ccnats"
	"github.com/sts-solutions/base-code/ccmsgqueue"
	"github.com/sts-solutions/base-code/ccmsgqueue/ccnats"
)

func (a *App) connect(ctx context.Context) (cnn Connectors, err error) {

	// Initialize tracer
	tracer, err := otel.InitTracer(
		a.ctx,
		"PoolService",
		config.App().Observability.Tracing.TraceHost,
		config.App().Observability.Tracing.TracePath,
		config.App().Env.Name,
		config.App().Observability.Tracing.SampleRate,
	)

	if err != nil {
		return cnn, errors.Wrap(err, "initializing tracer")
	}

	log.Logger().Info(ctx, "connecting to nats server")

	var natsConn ccmsgqueue.Connection = nil
	natsConn, err = ccnats.NewConnectionBuilder().
		WithHost(config.Infra().Nats.Host).
		WithPort(config.Infra().Nats.Port).
		WithMetrics(ccmetricsnats.NewConsumerMetrics(metrics.Namespace)).
		WithLogger(log.Logger()).
		WithTracer(tracer).
		Build()
	if err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("fail building nats connection, %v", err))
	}

	err = natsConn.Connect()
	if err != nil {
		log.Logger().Error(ctx, fmt.Sprintf("fail connecting to the message queue, %v", err))
	}

	go monitoringMsgQueueConnection(ctx, natsConn, config.App().Nats.ReconnectWait)

	log.Logger().Info(ctx, "connecting to mongo db")

	db, err := SetupMongoDB(a.ctx)
	if err != nil {
		return cnn, errors.Wrap(err, "setting up db")
	}

	cnn = Connectors{
		tracer: tracer,
		nats:   natsConn,
		db:     db,
		closeFunc: func() {
			if natsConn != nil {
				natsConn.Close()
			}
		},
	}

	return cnn, nil
}

func monitoringMsgQueueConnection(ctx context.Context, conn ccmsgqueue.Connection, wait time.Duration) {
	ticker := time.NewTicker(wait)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if !conn.IsConnected() {
				log.Logger().Info(ctx, "message queue disconnected, attempting to connect...")
				if err := conn.Connect(); err != nil || !conn.IsConnected() {
					log.Logger().Error(ctx, fmt.Sprintf("fail to reconnected to message queue, %v", err))
				} else {
					log.Logger().Info(ctx, "successfully reconnected to message queue")
				}
			}
		case <-ctx.Done():
			return
		}
	}
}
