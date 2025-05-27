package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/ronnyp07/SportStream/api/internal/app/httpserver"
	portsMetrics "github.com/ronnyp07/SportStream/api/internal/domain/ports/metrics"
	services "github.com/ronnyp07/SportStream/api/internal/domain/services/articles"
	"github.com/ronnyp07/SportStream/api/internal/metrics"
	"github.com/ronnyp07/SportStream/api/internal/pkg/config"
	"github.com/ronnyp07/SportStream/api/internal/pkg/infaestructure/database/repositories"
	"github.com/ronnyp07/SportStream/api/internal/pkg/infaestructure/log"

	"go.opentelemetry.io/otel/trace"
)

type Connectors struct {
	tracer    trace.Tracer
	closeFunc func()
	db        *MongoDB
}

type App struct {
	host        string
	ctx         context.Context
	ctxCancelFn func()
	connectors  Connectors
	termChan    chan os.Signal
	httpServer  *httpserver.Server
}

func (a *App) Start(ctx context.Context) error {
	a.ctx, a.ctxCancelFn = context.WithCancel(ctx)

	defer log.Logger().Info(a.ctx, "Application stopped successfully")

	var err error
	a.connectors, err = a.connect(ctx)
	if err != nil {
		return err
	}

	metricsHandler := metrics.NewMetricsHandler()
	metricsHandler.RegisterMetrics()

	appServices := setupServices(a.connectors, metricsHandler)

	server := httpserver.NewServerBuilder(httpserver.Services{
		ArticleService: appServices.ArticleServ,
	}).
		WithAddr(config.App().Http.HostAddress).
		WithReadTimeout(config.App().Http.ReadTimeout).
		WithWriteTimeout(config.App().Http.WriteTimeout).
		Build()

	server.Start(a.ctx, a.ctxCancelFn)
	a.httpServer = server

	a.termChan = make(chan os.Signal, 1)

	log.Logger().Info(a.ctx, fmt.Sprintf("Application services started successfully in %s", config.App().Env.Name))

	a.Shutdown()
	return nil
}

func (a *App) Shutdown() {
	signal.Notify(a.termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)

	log.Logger().Info(a.ctx, "stopping the server...")

	defer a.ctxCancelFn()
	defer a.connectors.closeFunc()

	defer func(httpSvr *httpserver.Server, ctx context.Context) {
		err := httpSvr.Stop(ctx)
		if err != nil {
			log.Logger().Fatal(a.ctx, "shutdown fail")
		}
	}(a.httpServer, a.ctx)

	var ctxErr error
	select {
	case <-a.termChan:
	case <-a.ctx.Done():
		ctxErr = a.ctx.Err()
		log.Logger().Error(a.ctx, fmt.Sprintf("context error, %v", ctxErr))
	}
}

func setupServices(c Connectors, metrics portsMetrics.MetricsHandler) Services {
	articleRepo := repositories.NewArticleRepository(c.db.DB, metrics)
	articlesServ := services.NewArticleService(articleRepo)
	// Implement service setup logic
	return Services{
		ArticleServ: articlesServ,
	}
}
