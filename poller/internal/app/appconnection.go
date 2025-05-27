package app

import (
	"context"
	"fmt"

	"emperror.dev/errors"
	"github.com/ronnyp07/SportStream/internal/pkg/config"
	"github.com/ronnyp07/SportStream/internal/pkg/infaestructure/log"
	appnats "github.com/ronnyp07/SportStream/internal/pkg/infaestructure/nats"
	"github.com/ronnyp07/SportStream/internal/pkg/otel"
	"github.com/sts-solutions/base-code/ccotel/ccotelnats"
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

	// Setup the database
	// db, err := SetupDB(a.ctx)
	// if err != nil {
	// 	return cnn, errors.Wrap(err, "setting up db")
	// }

	log.Logger().Info(ctx, fmt.Sprintf("connecting to nats server, %v", map[string]interface{}{
		"host": config.Infra().Nats.Host,
		"port": config.Infra().Nats.Port,
	}))

	nc, err := appnats.Connect(ctx,
		config.Infra().Nats.Host,
		config.Infra().Nats.Port)
	if err != nil {
		return cnn,
			errors.Wrap(err, "connecting to nats server")
	}

	js, err := nc.JetStream()
	if err != nil {
		return cnn, errors.Wrap(err, "connecting to jetstrem")
	}

	ccotelnats.SetTracer(nc, tracer)

	cnn = Connectors{
		tracer:    tracer,
		natsJSCtx: js,
		nats:      nc,
		closeFunc: func() {
			if config.Infra().MessageQueue.Enabled {
				nc.Close()
			}
		},
	}

	return cnn, nil
}
