package app

import (
	"context"

	"emperror.dev/errors"
	"github.com/ronnyp07/SportStream/api/internal/pkg/config"
	"github.com/ronnyp07/SportStream/api/internal/pkg/otel"
)

func (a *App) connect(ctx context.Context) (cnn Connectors, err error) {

	// Initialize tracer
	tracer, err := otel.InitTracer(
		a.ctx,
		"ApiService",
		config.App().Observability.Tracing.TraceHost,
		config.App().Observability.Tracing.TracePath,
		config.App().Env.Name,
		config.App().Observability.Tracing.SampleRate,
	)

	if err != nil {
		return cnn, errors.Wrap(err, "initializing tracer")
	}

	db, err := SetupMongoDB(a.ctx)
	if err != nil {
		return cnn, errors.Wrap(err, "setting up db")
	}

	cnn = Connectors{
		tracer: tracer,
		db:     db,
		closeFunc: func() {

		},
	}

	return cnn, nil
}
