package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/nats-io/nats.go"
	portMetrics "github.com/ronnyp07/SportStream/internal/domain/ports/metrics"
	"github.com/ronnyp07/SportStream/internal/domain/ports/services"
	"github.com/ronnyp07/SportStream/internal/domain/services/scheduler"
	"github.com/ronnyp07/SportStream/internal/metrics"
	"github.com/ronnyp07/SportStream/internal/pkg/config"
	"github.com/ronnyp07/SportStream/internal/pkg/infaestructure/log"
	natsQueue "github.com/ronnyp07/SportStream/internal/pkg/infaestructure/msgqueue"
	"go.opentelemetry.io/otel/trace"
)

type Connectors struct {
	nats      *nats.Conn
	natsJSCtx nats.JetStreamContext
	tracer    trace.Tracer
	closeFunc func()
}

type App struct {
	host              string
	ctx               context.Context
	ctxCancelFn       func()
	connectors        Connectors
	termChan          chan os.Signal
	shcedulerServ     scheduler.Service
	msgQueueProcessor services.MessageQueueProcessor
}

type connectors struct {
	msgQueue interface{} // Replace with actual message queue connector type
	// Add other connector fields as needed
}

func (a *App) Start(ctx context.Context) error {
	a.ctx, a.ctxCancelFn = context.WithCancel(ctx)

	defer log.Logger().Info(a.ctx, "Application stopped successfully")
	//log.Logger().Info(a.ctx, fmt.Sprintf("Application running in env: %s", "PoolService"))

	log.Logger().Info(a.ctx, fmt.Sprintf("Application services started successfully in %s, %v", config.App().Env.Name, config.App().Jobs))

	var err error
	a.connectors, err = a.connect(ctx)
	if err != nil {
		return err
	}

	metricsHandler := metrics.NewSchedulerMetricsHandler()
	metricsHandler.RegisterMetrics()

	appServices := setupServices(a.connectors)

	a.startScheduler(ctx, metricsHandler, appServices.MsgQueueService)

	//appServices := setupServices(a.connectors)
	appServices.MsgQueueService.PublishMessage(ctx, "SPORTSTREAM.DOCKER.status.updated", []byte("test nats"))
	//a.termChan = make(chan os.Signal, 1)
	//signal.Notify(a.termChan, os.Interrupt)

	// Create and start gRPC server
	// a.grpcServer = agrpc.NewAPI(a.ctx, appServices, config.Infra().Grpc.Port)
	// go a.grpcServer.Start(a.ctx, a.termChan)

	// Create and stxart HTTP API
	// a.httpAPI = http.NewAPI(appServices)
	// go a.httpAPI.Start(a.ctx, a.termChan)

	log.Logger().Info(a.ctx, fmt.Sprintf("Application services started successfully in %s, %v", config.App().Env.Name, config.App().Jobs))

	select {}

	// Wait for termination signal
	// <-a.termChan
	// a.Shutdown()
	//return nil
}

// func (a *app) connect() (connectors, error) {
// 	// Implement connection logic here
// 	return connectors{}, nil
// }

func (a *App) Shutdown() {
	signal.Notify(a.termChan, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-a.termChan

	log.Logger().Info(a.ctx, "stopping the server...")
	if err := a.shcedulerServ.Shutdown(); err != nil {
		log.Logger().Error(a.ctx, err.Error())
	}

	log.Logger().Info(a.ctx, "canceling main context")
	if a.ctxCancelFn != nil {
		a.ctxCancelFn()
	}
}

func (a *App) startScheduler(ctx context.Context,
	mectrics portMetrics.SchedulerMetricsHandler,
	msgService natsQueue.MsgQueueService) {
	schedulerService, err := scheduler.NewService(config.App().Jobs, mectrics, msgService)
	if err != nil {
		log.Logger().Info(ctx, "unable to start scheduler")
	}
	go schedulerService.Start(ctx)
	a.shcedulerServ = *schedulerService
}

func setupServices(c Connectors) Services {
	msgQueue := natsQueue.NewNatsService(c.natsJSCtx)
	// Implement service setup logic
	return Services{
		MsgQueueService: msgQueue,
	}
}
