package ports

import (
	"context"

	"github.com/go-co-op/gocron/v2"
	"github.com/ronnyp07/SportStream/internal/pkg/config"
)

type JobBuilder interface {
	BuildJob(ctx context.Context, name string, jobConfig config.Job, task func(ctx context.Context)) (gocron.Job, error)
}
