package ports

import (
	"context"

	"github.com/ronnyp07/SportStream/internal/pkg/config"
)

type IJob interface {
	Configure(ctx context.Context, jobCnfig config.Job) error

	Name() string
}
