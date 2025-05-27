package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"emperror.dev/errors"
	portsMetrics "github.com/ronnyp07/SportStream/worker/internal/domain/ports/metrics"
	"github.com/ronnyp07/SportStream/worker/internal/domain/ports/services"
	"github.com/ronnyp07/SportStream/worker/internal/domain/services/articles"
	"github.com/ronnyp07/SportStream/worker/internal/metrics"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/config"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/infaestructure/database/repositories"
	"github.com/ronnyp07/SportStream/worker/internal/pkg/infaestructure/log"
	subcriptions "github.com/ronnyp07/SportStream/worker/internal/pkg/infaestructure/natsconsumer"

	//ccmetricsnats "github.com/sts-solutions/base-code/ccmetrics/ccmsgqueue/ccnats"
	"github.com/sts-solutions/base-code/ccmsgqueue"
	"github.com/sts-solutions/base-code/ccmsgqueue/ccnats"
	"go.opentelemetry.io/otel/trace"
)

type Connectors struct {
	nats      ccmsgqueue.Connection
	tracer    trace.Tracer
	closeFunc func()
	db        *MongoDB
}

type App struct {
	host              string
	ctx               context.Context
	ctxCancelFn       func()
	connectors        Connectors
	termChan          chan os.Signal
	msgQueueProcessor services.MessageQueueProcessor
}

func (a *App) Start(ctx context.Context) error {
	a.ctx, a.ctxCancelFn = context.WithCancel(ctx)

	defer log.Logger().Info(a.ctx, "Application stopped successfully")
	//log.Logger().Info(a.ctx, fmt.Sprintf("Application running in env: %s", "PoolService"))

	var err error
	a.connectors, err = a.connect(ctx)
	if err != nil {
		return err
	}

	metricsHandler := metrics.NewMetricsHandler()
	metricsHandler.RegisterMetrics()

	appServices := setupServices(a.connectors, metricsHandler)

	natsServices := subcriptions.Service{
		ArticleServ: appServices.ArticleServ,
	}

	natsMessageHandler := subcriptions.NewHandler(&natsServices)
	natsConsumer, err := ccnats.NewConsumerBuilder().
		WithName(config.App().Nats.Consumers.Articles.Update.ConsumerName).
		WithStream(config.App().Nats.Consumers.Articles.Update.Stream).
		WithSubject(config.App().Nats.Consumers.Articles.Update.Subject).
		WithConnection(a.connectors.nats).
		//WithMetrics(ccmetricsnats.NewConsumerMetrics(metrics.Namespace)).
		WithMessageHandler(natsMessageHandler.HandleMessage).
		Build()

	if err != nil {
		return errors.Wrap(err, "building nats consumer")
	}

	go func() {
		err = natsConsumer.Consume(a.ctx)
		if err != nil {
			log.Logger().Error(a.ctx, fmt.Sprintf("nats consumer error, %v", err))
			os.Exit(1)
		}
	}()

	//appServices := setupServices(a.connectors)
	//a.termChan = make(chan os.Signal, 1)
	//signal.Notify(a.termChan, os.Interrupt)

	// Create and start gRPC server
	// a.grpcServer = agrpc.NewAPI(a.ctx, appServices, config.Infra().Grpc.Port)
	// go a.grpcServer.Start(a.ctx, a.termChan)

	// Create and stxart HTTP API
	// a.httpAPI = http.NewAPI(appServices)
	// go a.httpAPI.Start(a.ctx, a.termChan)

	log.Logger().Info(a.ctx, fmt.Sprintf("Application services started successfully in %s", config.App().Env.Name))

	select {}

	// Wait for termination signal
	// <-a.termChan
	// a.Shutdown()
	//return nil
}

func (a *App) Shutdown() {
	signal.Notify(a.termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-a.termChan

	log.Logger().Info(a.ctx, "stopping the server...")

	log.Logger().Info(a.ctx, "canceling main context")
	if a.ctxCancelFn != nil {
		a.ctxCancelFn()
	}
}

func setupServices(c Connectors, metrics portsMetrics.MetricsHandler) Services {
	articleRepo := repositories.NewArticleRepository(c.db.DB, metrics)
	articlesServ := articles.NewArticlesService(articleRepo)
	// Implement service setup logic
	return Services{
		ArticleServ: articlesServ,
	}
}
