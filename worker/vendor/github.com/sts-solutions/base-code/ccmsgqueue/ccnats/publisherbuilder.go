package ccnats

import (
	"errors"

	ccmetrics "github.com/sts-solutions/base-code/ccmetrics/ccmsgqueue"
	"github.com/sts-solutions/base-code/ccmsgqueue"
	"github.com/sts-solutions/base-code/ccvalidation"
)

// PublisherBuilder is a builder for constructing a publisher.
type PublisherBuilder struct {
	publisher *publisher
	metrics   ccmetrics.PublisherMetrics
}

// NewPublisherBuilder creates a new instance of PublisherBuilder.
func NewPublisherBuilder() *PublisherBuilder {
	return &PublisherBuilder{
		publisher: &publisher{},
	}
}

// WithConnection sets the connection for the publisher. Connection is mandatory.
func (pb *PublisherBuilder) WithConnection(conn ccmsgqueue.Connection) *PublisherBuilder {
	if c, ok := conn.(*connection); ok {
		pb.publisher.conn = c
	}
	return pb
}

// WithMetrics sets the metrics for the publisher.
func (pb *PublisherBuilder) WithMetrics(metrics ccmetrics.PublisherMetrics) *PublisherBuilder {
	pb.metrics = metrics
	return pb
}

// Build validates the configuration and returns the constructed Publisher.
func (pb *PublisherBuilder) Build() (ccmsgqueue.Publisher, error) {
	result := pb.validate()
	if result.IsFailure() {
		return nil, errors.New(result.Error())
	}
	pb.publisher.metrics = ccmsgqueue.NewPublisherMetrics(pb.metrics)
	return pb.publisher, nil
}

// validate checks if all required fields are set.
func (pb *PublisherBuilder) validate() ccvalidation.Result {
	result := ccvalidation.Result{}

	if pb.publisher.conn == nil {
		result.AddErrorMessage("Publisher connection is missing")
	}
	return result
}
