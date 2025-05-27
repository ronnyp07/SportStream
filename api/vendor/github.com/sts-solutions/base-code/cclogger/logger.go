package cclogger

import (
	"context"

	"github.com/sts-solutions/base-code/cclog"
	"github.com/sts-solutions/base-code/ccmiddlewares/cccorrelation"
	"go.uber.org/zap/zapcore"
)

type Logger interface {
	WithField(name string, value interface{}) Logger
	Debug(ctx context.Context, msg string, logFields ...LogField)
	Info(ctx context.Context, msg string, logFields ...LogField)
	Warn(ctx context.Context, msg string, logFields ...LogField)
	Error(ctx context.Context, msg string, logFields ...LogField)
	Fatal(ctx context.Context, msg string, logFields ...LogField)
}

type logger struct {
	addCorrelationID bool
	Logger           cclog.Logger // embedded interface
	fields           []LogField
}

func (t *logger) WithField(name string, value interface{}) Logger {
	t.fields = append(t.fields, LogField{Key: name, Value: value})
	return t
}

func (t *logger) Debug(ctx context.Context, msg string, logFields ...LogField) {
	t.Logger.Debug(ctx, msg, t.getFields(ctx, logFields...)...)
	t.fields = []LogField{}
}

func (t *logger) Info(ctx context.Context, msg string, logFields ...LogField) {
	t.Logger.Info(ctx, msg, t.getFields(ctx, logFields...)...)
	t.fields = []LogField{}
}

func (t *logger) Warn(ctx context.Context, msg string, logFields ...LogField) {
	t.Logger.Warn(ctx, msg, t.getFields(ctx, logFields...)...)
	t.fields = []LogField{}
}

func (t *logger) Error(ctx context.Context, msg string, logFields ...LogField) {
	t.Logger.Error(ctx, msg, t.getFields(ctx, logFields...)...)
	t.fields = []LogField{}
}

func (t *logger) Fatal(ctx context.Context, msg string, logFields ...LogField) {
	t.Logger.Fatal(ctx, msg, t.getFields(ctx, logFields...)...)
	t.fields = []LogField{}
}

func (t *logger) getFields(ctx context.Context, logFields ...LogField) []zapcore.Field {
	var fields []zapcore.Field

	if t.addCorrelationID {
		logFields = addCorrelationID(ctx, logFields)
	}

	logFields = append(logFields, t.fields...)
	fields = toZapCoreFields(logFields...)
	return fields
}

func addCorrelationID(ctx context.Context, logFields []LogField) []LogField {
	corrID := cccorrelation.GetCorrelationID(ctx)
	logFields = append(logFields, LogField{Key: cccorrelation.LogKey, Value: corrID})
	return logFields
}
