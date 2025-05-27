package ccnats

import (
	"github.com/sts-solutions/base-code/cclogger"
	ccmetrics "github.com/sts-solutions/base-code/ccmetrics/ccmsgqueue"
	"github.com/sts-solutions/base-code/ccmsgqueue"
	"github.com/sts-solutions/base-code/ccvalidation"

	"emperror.dev/errors"
	"go.opentelemetry.io/otel/trace"
)

// ConnectionBuilder is a builder for constructing a Connection.
type ConnectionBuilder struct {
	connection *connection
	metrics    ccmetrics.ConnectionMetrics
	tracer     trace.Tracer
	logger     cclogger.Logger
}

// NewConnectionBuilder creates a new instance of ConnectionBuilder.
func NewConnectionBuilder() *ConnectionBuilder {
	return &ConnectionBuilder{
		connection: &connection{},
	}
}

// WithPort sets the port for the connection. Port is mandatory.
func (cb *ConnectionBuilder) WithPort(port int) *ConnectionBuilder {
	cb.connection.port = port
	return cb
}

// WithHost sets the host for the connection. Host is mandatory.
func (cb *ConnectionBuilder) WithHost(host string) *ConnectionBuilder {
	cb.connection.host = host
	return cb
}

// WithMetrics sets the metrics for the connection.
func (cb *ConnectionBuilder) WithMetrics(metrics ccmetrics.ConnectionMetrics) *ConnectionBuilder {
	cb.metrics = metrics
	return cb
}

// WithTracer sets the tracer for the connection.
func (cb *ConnectionBuilder) WithTracer(tracer trace.Tracer) *ConnectionBuilder {
	cb.tracer = tracer
	return cb
}

// WithLogger sets the logger for the connection.
func (cb *ConnectionBuilder) WithLogger(logger cclogger.Logger) *ConnectionBuilder {
	cb.logger = logger
	return cb
}

// Build validates the configuration and returns the constructed Connection.
func (cb *ConnectionBuilder) Build() (ccmsgqueue.Connection, error) {
	result := cb.validate()
	if result.IsFailure() {
		return nil, errors.Wrap(result, "validating nats connection builder")
	}

	cb.connection.metricts = ccmsgqueue.NewConnectionMetrics(cb.metrics)
	cb.connection.tracer = NewTracer(cb.tracer)
	cb.connection.logger = ccmsgqueue.NewLogger(cb.logger)

	return cb.connection, nil
}

// validate checks if all required fields are set.
func (cb *ConnectionBuilder) validate() ccvalidation.Result {
	result := ccvalidation.Result{}

	if cb.connection.port == 0 {
		result.AddErrorMessage("Connection port is missing")
	}

	if cb.connection.host == "" {
		result.AddErrorMessage("Connection host is missing")
	}

	return result
}
