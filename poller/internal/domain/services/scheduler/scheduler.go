package scheduler

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/ronnyp07/SportStream/internal/domain/models"
	ports "github.com/ronnyp07/SportStream/internal/domain/ports/job"
	"github.com/ronnyp07/SportStream/internal/domain/ports/metrics"
	"github.com/ronnyp07/SportStream/internal/domain/services/jobs/jobbuilder"
	pooller "github.com/ronnyp07/SportStream/internal/domain/services/jobs/poller"
	"github.com/ronnyp07/SportStream/internal/pkg/config"
	"github.com/ronnyp07/SportStream/internal/pkg/infaestructure/log"
	natsQueue "github.com/ronnyp07/SportStream/internal/pkg/infaestructure/msgqueue"
	"github.com/sts-solutions/base-code/cchttp"
)

type Service struct {
	scheduler      gocron.Scheduler
	scheduledJobs  []ports.IJob
	jobsConfig     config.Jobs
	metricsHandler metrics.SchedulerMetricsHandler
	msgQueueServ   natsQueue.MsgQueueService
}

func NewService(jobsConfig config.Jobs,
	metrics metrics.SchedulerMetricsHandler,
	msgQueueServ natsQueue.MsgQueueService) (*Service, error) {
	sch, err := gocron.NewScheduler()
	if err != nil {
		log.Logger().Fatal(context.Background(), fmt.Sprintf("Cannot create go cron scheduler due to %w", err))
		return nil, err
	}

	jobBuilder := jobbuilder.New(sch)
	httpclient := cchttp.NewClient(0, 0, 0, time.Minute)

	poollerJob := pooller.New(sch, jobBuilder, httpclient, metrics, msgQueueServ)

	return &Service{
		scheduler: sch,
		scheduledJobs: []ports.IJob{
			poollerJob,
		},
		jobsConfig:     jobsConfig,
		metricsHandler: metrics,
	}, nil
}

func (s *Service) Start(ctx context.Context) {
	log.Logger().Info(ctx, "Configuring jobs")
	s.configureJobs(ctx)

	log.Logger().Info(ctx, "Starting cronjob scheduler")
	s.scheduler.Start()
}

func (s *Service) Shutdown() error {
	return s.scheduler.Shutdown()
}

func (s *Service) GetAllJobs(ctx context.Context) ([]models.Job, error) {
	log.Logger().Info(ctx, "Preparing to get all jobs")
	jobs := make([]models.Job, 0)
	for _, job := range s.scheduler.Jobs() {
		jobs = append(jobs, fromLibraryToDomain(job))
	}
	if len(jobs) == 0 {
		return nil, errors.New("job not found")
	}
	log.Logger().Info(ctx, fmt.Sprintf("Getting all jobs %v", jobs))
	return jobs, nil
}

func (s *Service) RunJobNow(ctx context.Context, ID string) (models.Job, error) {
	log.Logger().Info(ctx, fmt.Sprintf("Run job now %v", ID))
	for _, scheduledJob := range s.scheduler.Jobs() {
		if scheduledJob.ID().String() == ID {
			log.Logger().Info(ctx, fmt.Sprintf("Run job now %v", ID))
			return fromLibraryToDomain(scheduledJob), scheduledJob.RunNow()
		}
	}
	return models.Job{}, errors.New("not found")
}

func (s *Service) configureJobs(ctx context.Context) {
	log.Logger().Info(ctx, "Configuring all jobs")
	for _, scheduledJob := range s.scheduledJobs {
		jobName := scheduledJob.Name()
		if err := scheduledJob.Configure(ctx, s.jobsConfig[jobName]); err != nil {
			log.Logger().Error(ctx, fmt.Sprintf("Error configuring job %s due to %v", jobName, err))
		}
	}
}

func fromLibraryToDomain(job gocron.Job) models.Job {
	nextRun, _ := job.NextRun()
	lastRun, _ := job.LastRun()
	return models.Job{
		ID:      job.ID().String(),
		Name:    job.Name(),
		NextRun: nextRun,
		LastRun: lastRun,
		Tags:    job.Tags(),
	}
}
