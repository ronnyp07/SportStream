package jobbuilder

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/go-co-op/gocron/v2"
	"github.com/ronnyp07/SportStream/internal/pkg/config"
	"github.com/ronnyp07/SportStream/internal/pkg/infaestructure/log"
)

const (
	CronJobType     = "CRONJOB"
	DurationJobType = "DURATIONJOB"
)

type JobBuilder struct {
	scheduler gocron.Scheduler
}

func New(scheduler gocron.Scheduler) *JobBuilder {
	return &JobBuilder{
		scheduler: scheduler,
	}
}

func (j *JobBuilder) BuildJob(ctx context.Context, name string, cfg config.Job, fn func(ctx context.Context)) (gocron.Job, error) {
	if !cfg.Enabled {
		log.Logger().Info(ctx, fmt.Sprintf("Job is disabled name %s", name))
		return nil, nil
	}

	jobDefinition, err := getJobDefinition(cfg.Type, cfg.Interval, cfg.UseSeconds)
	if err != nil {
		return nil, err
	}

	return j.scheduler.NewJob(
		jobDefinition,
		gocron.NewTask(fn, ctx),
		gocron.WithName(name),
	)
}

func getJobDefinition(jType, interval string, withSeconds bool) (gocron.JobDefinition, error) {
	switch jType {
	case CronJobType:
		return gocron.CronJob(interval, withSeconds), nil
	case DurationJobType:
		duration, err := time.ParseDuration(interval)
		if err != nil {
			return nil, err
		}
		return gocron.DurationJob(duration), nil
	default:
		return nil, errors.New("invalid job type")
	}
}
