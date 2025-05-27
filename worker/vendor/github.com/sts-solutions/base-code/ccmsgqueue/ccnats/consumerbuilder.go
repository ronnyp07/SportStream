package ccnats

import (
	"context"

	ccmetrics "github.com/sts-solutions/base-code/ccmetrics/ccmsgqueue"
	"github.com/sts-solutions/base-code/ccmsgqueue"
	"github.com/sts-solutions/base-code/ccvalidation"

	"emperror.dev/errors"
)

// ConsumerBuilder is a builder for constructing a NATS consumer
type ConsumerBuilder struct {
	consumer consumer
	metrics  ccmetrics.ConsumerMetrics
}

// NewConsumerBuilder creates a new instance of ConsumerBuilder
func NewConsumerBuilder() *ConsumerBuilder {
	return &ConsumerBuilder{
		consumer: consumer{},
	}
}

// WithNamespace sets the namespace for the consumer
func (cb *ConsumerBuilder) WithNamespace(namespace string) *ConsumerBuilder {
	cb.consumer.namespace = namespace
	return cb
}

// WithStream sets the stream for the consumer
func (cb *ConsumerBuilder) WithStream(stream string) *ConsumerBuilder {
	cb.consumer.stream = stream
	return cb
}

// WithSubject sets the subject for the consumer
func (cb *ConsumerBuilder) WithSubject(subject string) *ConsumerBuilder {
	cb.consumer.subject = subject
	return cb
}

// WithName sets the name for the consumer
func (cb *ConsumerBuilder) WithName(name string) *ConsumerBuilder {
	cb.consumer.name = name
	return cb
}

// WithConnection sets the connection for the consumer
func (cb *ConsumerBuilder) WithConnection(conn ccmsgqueue.Connection) *ConsumerBuilder {
	if c, ok := conn.(*connection); ok {
		cb.consumer.conn = c
	}
	return cb
}

// WithMessageHandler sets the message handler for the consumer
func (cb *ConsumerBuilder) WithMessageHandler(handler func(ctx context.Context, msg ccmsgqueue.ConsumeMessage)) *ConsumerBuilder {
	cb.consumer.msgHandler = handler
	return cb
}

// WithMetrics sets the metrics for the consumer
func (cb *ConsumerBuilder) WithMetrics(metrics ccmetrics.ConsumerMetrics) *ConsumerBuilder {
	cb.metrics = metrics
	return cb
}

// Build validates the configuration and returns the constructed consumer
func (cb *ConsumerBuilder) Build() (*consumer, error) {
	result := cb.validate()
	if result.IsFailure() {
		return nil, errors.Wrap(result, "validating nats consuemr builder")
	}

	cb.consumer.metrics = ccmsgqueue.NewConsumerMetrics(cb.metrics)
	return &cb.consumer, nil
}

// validate checks if all required fields are set
func (cb *ConsumerBuilder) validate() ccvalidation.Result {
	result := ccvalidation.Result{}

	if cb.consumer.stream == "" {
		result.AddErrorMessage("consumer stream is missing")
	}
	if cb.consumer.subject == "" {
		result.AddErrorMessage("consumer subject is missing")
	}
	if cb.consumer.name == "" {
		result.AddErrorMessage("consumer name is missing")
	}
	if cb.consumer.msgHandler == nil {
		result.AddErrorMessage("consumer message handler is missing")
	}

	return result
}
