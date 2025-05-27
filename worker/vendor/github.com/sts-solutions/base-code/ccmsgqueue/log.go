package ccmsgqueue

import (
	"context"

	"github.com/sts-solutions/base-code/cclogger"
)

type Logger interface {
	LogInfo(ctx context.Context, msg string, logFields ...cclogger.LogField)
	LogWarn(ctx context.Context, msg string, logFields ...cclogger.LogField)
	LogError(ctx context.Context, msg string, logFields ...cclogger.LogField)
}

type logger struct {
	logger    cclogger.Logger
	shouldLog bool
}

func NewLogger(l cclogger.Logger) Logger {
	return &logger{
		logger:    l,
		shouldLog: l != nil,
	}
}

func (l logger) LogInfo(ctx context.Context, msg string, logFields ...cclogger.LogField) {
	if l.shouldLog {
		l.logger.Info(ctx, msg, logFields...)
	}
}

func (l logger) LogWarn(ctx context.Context, msg string, logFields ...cclogger.LogField) {
	if l.shouldLog {
		l.logger.Warn(ctx, msg, logFields...)
	}
}

func (l logger) LogError(ctx context.Context, msg string, logFields ...cclogger.LogField) {
	if l.shouldLog {
		l.logger.Error(ctx, msg, logFields...)
	}
}
