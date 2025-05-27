package cclog

import (
	"context"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// LoggerBuilder is a builder for creating a Logger.
type LoggerBuilder struct {
	// The log Level. Default is Debug.
	level Level
	// The number of stack frames to skip when logging. Default is 0.
	callerskip int
	// Fields to be added to every log entry. Default is empty.
	resourceFields []zapcore.Field
	// Write syncer to write logs to. Default is os.Stdout.
	ws zapcore.WriteSyncer
	// A function to extract attributes from the context. Default is nil.
	attributeExtractor func(context.Context) []zapcore.Field
}

// NewBuilder creates a new LoggerBuilder.
func NewBuilder() LoggerBuilder {
	return LoggerBuilder{
		level:          Debug,
		callerskip:     0,
		resourceFields: []zapcore.Field{},
		ws:             zapcore.AddSync(os.Stdout),
	}
}

// WithLevel sets the log level.
func (b LoggerBuilder) WithLevel(level Level) LoggerBuilder {
	b.level = level
	return b
}

// WithCallerskip sets the number of stack frames to skip when logging.
func (b LoggerBuilder) WithCallerskip(skip int) LoggerBuilder {
	b.callerskip = skip
	return b
}

// WithResourceFields sets the fields to be added to every log entry.
func (b LoggerBuilder) WithResourceFields(fields ...zapcore.Field) LoggerBuilder {
	b.resourceFields = fields
	return b
}

// WithOutput sets write syncer to write logs to.
func (b LoggerBuilder) WithOutput(ws zapcore.WriteSyncer) LoggerBuilder {
	b.ws = ws
	return b
}

// WithAttributeExtractor sets the function to extract attributes from the context.
func (b LoggerBuilder) WithAttributeExtractor(extractor func(context.Context) []zapcore.Field) LoggerBuilder {
	b.attributeExtractor = extractor
	return b
}

// Build creates a Logger.
func (b LoggerBuilder) Build() Logger {
	encoder := zapcore.NewJSONEncoder(defaultEncoderConfig())
	writer := b.ws

	var enabler zapcore.LevelEnabler
	switch b.level {
	case Debug:
		enabler = zapcore.DebugLevel
	case Info:
		enabler = zapcore.InfoLevel
	case Warn:
		enabler = zapcore.WarnLevel
	case Error:
		enabler = zapcore.ErrorLevel
	}

	loggerCore := zapcore.NewCore(encoder, writer, enabler)
	zapLogger := zap.New(loggerCore, zap.AddCallerSkip(b.callerskip), zap.AddCaller())

	return Logger{
		logger: zapLogger,
		resources: &Attributes{
			Fields: b.resourceFields,
		},
		attrExtractor: b.attributeExtractor,
	}
}
