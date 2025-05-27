package cclog

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const (
	// ATTRIBUTES_NAME constants
	AttributesName    = "Attributes"
	Body              = "Body"
	ObservedTimestamp = "ObservedTimestamp"
	Resources         = "Resources"
	SeverityNumber    = "SeverityNumber"
	SeverityText      = "SeverityText"
	SpanId            = "SpanId"
	SourceType        = "SourceType"
	Timestamp         = "Timestamp"
	TraceFlags        = "TraceFlags"
	TraceId           = "TraceId"

	// OBSERVED_TIMESTAMP constants
	ResourcesName       = "Resources"
	SeverityNumberConst = "SeverityNumber"
	SeverityTextConst   = "SeverityText"
	SpanIdName          = "SpanId"
	SourceTypeConst     = "SourceType"
	TimestampConst      = "Timestamp"
	TraceFlagsConst     = "TraceFlags"
	TraceIdName         = "TraceId"
)

type ContextKey string
type Level uint8

const (
	ATTRIBUTES_KEY = ContextKey("_sts")
)

const (
	Debug Level = iota
	Info
	Warn
	Error
)

// Logger is the instance of loggers which is a wrapper on top of Uber Zap.
type Logger struct {
	// The underlying zap logger instance
	logger *zap.Logger

	// The resource part of the otel log record
	resources *Attributes

	// Optional extractor for overwriting attributes given by the context Attribute value
	attrExtractor func(context.Context) []zapcore.Field
}

// NewLogger instantiates a logger and extracts span id, trace id and attributes from the context
// or uses the specified values
// Deprecated: Use LoggerBuilder instead.
func NewLogger(lvl Level, resourceFields ...zapcore.Field) Logger {
	// Set the caller skip to 1, because we are wrapping zap logger with one more func call
	logger := newZapCoreLogger(lvl, 1)
	return Logger{
		logger: logger,
		resources: &Attributes{
			Fields: resourceFields,
		},
	}
}

// NewLoggerWithCallerSkip instantiates a logger with a specified caller skip and the given
// resource fields.
// Deprecated: Use LoggerBuilder instead.
func NewLoggerWithCallerSkip(lvl Level, callerSkip int, resourceFields ...zapcore.Field) Logger {
	// If the caller passes 1 as a caller skip, it really means 1 + 1, because we are wrapping zap
	// logger logging funcs.
	logger := newZapCoreLogger(lvl, callerSkip+1)
	return Logger{
		logger: logger,
		resources: &Attributes{
			Fields: resourceFields,
		},
	}
}

// NewLoggerInternal creates a logger with an existing zap.Logger instance
func NewLoggerInternal(logger *zap.Logger, resourceFields ...zapcore.Field) Logger {
	return Logger{
		logger: logger,
		resources: &Attributes{
			Fields: resourceFields,
		},
	}
}

// NewTestLogger helps the test to capture stdout.
// Deprecated: Use LoggerBuilder instead.
func NewTestLogger(resourceFields ...zapcore.Field) (Logger, *CaptureWriter) {
	cw := CaptureWriter{}
	encoder := zapcore.NewJSONEncoder(defaultEncoderConfig())
	writer := &cw
	enabler := zapcore.InfoLevel
	loggerCore := zapcore.NewCore(encoder, writer, enabler)
	zapLogger := zap.New(loggerCore, zap.AddCallerSkip(1), zap.AddCaller())
	logger := NewLoggerInternal(zapLogger, resourceFields...)
	return logger, writer
}

// WithAttributeExtractor sets up an attribute extractor for the logger.
// Attribute extractor can get values from the context.
func (l Logger) WithAttributeExtractor(extractor func(context.Context) []zapcore.Field) {
	l.attrExtractor = extractor
}

// Debug logs a message at debug level, adding the fields to the Attributes map of the telemetry message.
func (l *Logger) Debug(ctx context.Context, msg string, attrFields ...zapcore.Field) {
	topFields := l.defaultFields(ctx)
	topFields = append(topFields, []zapcore.Field{
		zap.String(SeverityText, "DEBUG"),
		zap.Int(SeverityNumber, 5),
		zap.String(Body, msg),
	}...)
	topFields = l.mergeAttributesToFieldList(ctx, topFields, attrFields...)
	l.logger.Debug(msg, topFields...)
}

// Info logs a message at info level, adding the fields to the Attributes map of the telemetry message.
func (l *Logger) Info(ctx context.Context, msg string, attrFields ...zapcore.Field) {
	topFields := l.defaultFields(ctx)
	topFields = append(topFields, []zapcore.Field{
		zap.String(SeverityText, "INFO"),
		zap.Int(SeverityNumber, 9),
		zap.String(Body, msg),
	}...)
	topFields = l.mergeAttributesToFieldList(ctx, topFields, attrFields...)
	l.logger.Info(msg, topFields...)
}

// Warn logs a message at warning level, adding the fields to the Attributes map of the telemetry message.
func (l *Logger) Warn(ctx context.Context, msg string, attrFields ...zapcore.Field) {
	topFields := l.defaultFields(ctx)
	topFields = append(topFields, []zapcore.Field{
		zap.String(SeverityText, "WARN"),
		zap.Int(SeverityNumber, 13),
		zap.String(Body, msg),
	}...)
	topFields = l.mergeAttributesToFieldList(ctx, topFields, attrFields...)
	l.logger.Warn(msg, topFields...)
}

// Error logs a message at error level, adding the fields to the Attributes map of the telemetry message.
func (l *Logger) Error(ctx context.Context, msg string, attrFields ...zapcore.Field) {
	topFields := l.defaultFields(ctx)
	topFields = append(topFields, []zapcore.Field{
		zap.String(SeverityText, "ERROR"),
		zap.Int(SeverityNumber, 17),
		zap.String(Body, msg),
	}...)
	topFields = l.mergeAttributesToFieldList(ctx, topFields, attrFields...)
	l.logger.Error(msg, topFields...)
}

// Fatal logs a message at fatal level, adding the fields to the Attributes map of the telemetry message.
func (l *Logger) Fatal(ctx context.Context, msg string, attrFields ...zapcore.Field) {
	topFields := l.defaultFields(ctx)
	topFields = append(topFields, []zapcore.Field{
		zap.String(SeverityText, "FATAL"),
		zap.Int(SeverityNumber, 21),
		zap.String(Body, msg),
	}...)
	topFields = l.mergeAttributesToFieldList(ctx, topFields, attrFields...)
	l.logger.Fatal(msg, topFields...)
}

// Sync flushes any buffered log entries.
func (l *Logger) Sync() error {
	return l.logger.Sync()
}

// Bool wraps a zapcore bool field. It can be used in collecting [Attributes] and passing
// parameters to [Logger.Debug], [Logger.Info] etc functions.
// Deprecated: Use NewField instead.
func Bool(key string, value bool) zapcore.Field {
	return zap.Bool(key, value)
}

// Int wraps a zapcore int field. It can be used in collecting [Attributes] and passing
// parameters to [Logger.Debug], [Logger.Info] etc functions.
// Deprecated: Use NewField instead.
func Int(key string, value int) zapcore.Field {
	return zap.Int(key, value)
}

// String wraps a zapcore string field. It can be used in collecting [Attributes] and passing
// parameters to [Logger.Debug], [Logger.Info] etc functions.
// Deprecated: Use NewField instead.
func String(key string, value string) zapcore.Field {
	return zap.String(key, value)
}

func NewField(key string, value any) zapcore.Field {
	switch v := value.(type) {
	case string:
		return zap.String(key, v)
	case int:
		return zap.Int(key, v)
	case int64:
		return zap.Int64(key, v)
	case float64:
		return zap.Float64(key, v)
	case bool:
		return zap.Bool(key, v)
	case time.Time:
		return zap.Time(key, v)
	case time.Duration:
		return zap.Duration(key, v)
	case error:
		return zap.Error(v)
	default:
		return zap.Any(key, v)
	}
}

func (l *Logger) defaultFields(ctx context.Context) []zapcore.Field {
	fields := make([]zapcore.Field, 0, 3)

	fields = append(fields, zap.Int(TraceFlags, 0))

	spanCtx := trace.SpanContextFromContext(ctx)
	if spanCtx.HasTraceID() {
		fields = append(fields, zap.String(TraceIdName, spanCtx.TraceID().String()))
	}

	if spanCtx.HasSpanID() {
		fields = append(fields, zap.String(SpanIdName, spanCtx.SpanID().String()))
	}

	fields = append(fields, zap.Object(ResourcesName, l.resources))

	return fields

}

func (l *Logger) mergeAttributesToFieldList(ctx context.Context, fields []zapcore.Field, paramAttrs ...zapcore.Field) []zapcore.Field {
	var attr *Attributes

	if attr = extractAttributesFromContext(ctx, ATTRIBUTES_KEY); attr == nil {
		attr = &Attributes{
			Fields: make([]zapcore.Field, 0),
		}
	}

	for i := range paramAttrs {
		attr.Fields = append(attr.Fields, paramAttrs[i])
	}

	return append(fields, zap.Object(AttributesName, attr))
}

func defaultEncoderConfig() zapcore.EncoderConfig {
	config := zap.NewProductionEncoderConfig()
	config.MessageKey = "message"
	return config
}

func newZapCoreLogger(lvl Level, callerSkip int) *zap.Logger {
	encoder := zapcore.NewJSONEncoder(defaultEncoderConfig())
	writer := zapcore.AddSync(os.Stdout)
	var enabler zapcore.LevelEnabler
	switch lvl {
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
	zapLogger := zap.New(loggerCore, zap.AddCallerSkip(callerSkip), zap.AddCaller())
	return zapLogger
}

func extractStringFromContext(ctx context.Context, contextKey ContextKey, fieldName string) zapcore.Field {
	value := ctx.Value(contextKey)
	if value != nil {
		// Convert value to string
		field := zap.String(fieldName, fmt.Sprintf("%v", value))
		return field
	}
	return zap.Skip()
}

func extractAttributesFromContext(ctx context.Context, contextKey ContextKey) *Attributes {
	if value := ctx.Value(contextKey); value != nil {
		if val, ok := value.(*Attributes); ok {
			return val
		}
	}
	return nil
}
