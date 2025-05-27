package cclogger

import (
	"strings"

	"github.com/sts-solutions/base-code/cclog"
)

const (
	defaultCallerSkip = 2
	appNameKey        = "app"
)

type LoggerBuilder struct {
	appName          string
	level            Level
	logFields        []LogField
	addCorrelationID bool
}

func NewBuilder() *LoggerBuilder {
	return &LoggerBuilder{
		level:            Debug,
		logFields:        []LogField{},
		addCorrelationID: false,
	}
}

func (b *LoggerBuilder) WithLevel(level Level) *LoggerBuilder {
	b.level = level
	return b
}

func (b *LoggerBuilder) WithStringLevel(level string) *LoggerBuilder {
	level = strings.ToLower(level)
	switch level {
	case "debug":
		b.level = Debug
	case "info":
		b.level = Info
	case "warn":
		b.level = Warn
	case "error":
		b.level = Error
	default:
		b.level = Debug
	}

	return b
}

func (b *LoggerBuilder) WithAppName(name string) *LoggerBuilder {
	b.appName = name
	return b
}

func (b *LoggerBuilder) WithFields(logFields ...LogField) *LoggerBuilder {
	b.logFields = logFields
	return b
}

func (b *LoggerBuilder) WithCorrelationID() *LoggerBuilder {
	b.addCorrelationID = true
	return b
}

func (b *LoggerBuilder) Build() (Logger, error) {
	b.logFields = append(b.logFields, []LogField{
		{
			Key:   appNameKey,
			Value: b.appName,
		},
	}...)

	logLvl := cclog.Level(b.level)
	loggerBuilder := cclog.NewBuilder().
		WithLevel(logLvl).
		WithResourceFields(toZapCoreFields(b.logFields...)...)

	return &logger{
		Logger:           loggerBuilder.Build(),
		addCorrelationID: b.addCorrelationID,
	}, nil

}
