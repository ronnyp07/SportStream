package ccmsgqueue

import ccmetrics "github.com/sts-solutions/base-code/ccmetrics/ccmsgqueue"

type ConnectionMetrics interface {
	ErrorInc(errorType string)
}

type connectionMetrics struct {
	metrics      ccmetrics.ConnectionMetrics
	shouldReport bool
}

func NewConnectionMetrics(metrics ccmetrics.ConnectionMetrics) ConnectionMetrics {
	return connectionMetrics{
		metrics:      metrics,
		shouldReport: metrics != nil,
	}
}

func (m connectionMetrics) ErrorInc(errorType string) {
	if m.shouldReport {
		m.metrics.ErrorInc(errorType)
	}
}

type ConsumerMetrics interface {
	MessageLag(duration float64, subject, consumerName string)
	MessageCounterInc(subject, consumerName string)
	ErrorInc(errorType string)
}

type consumerMetrics struct {
	metrics      ccmetrics.ConsumerMetrics
	shouldReport bool
}

func NewConsumerMetrics(metrics ccmetrics.ConsumerMetrics) ConsumerMetrics {
	return consumerMetrics{
		metrics:      metrics,
		shouldReport: metrics != nil,
	}
}

func (m consumerMetrics) MessageLag(duration float64, subject, consumerName string) {
	if m.shouldReport {
		m.metrics.MessageLag(duration, subject, consumerName)
	}
}

func (m consumerMetrics) MessageCounterInc(subject, consumerName string) {
	if m.shouldReport {
		m.metrics.MessageCounterInc(subject, consumerName)
	}
}

func (m consumerMetrics) ErrorInc(errorType string) {
	if m.shouldReport {
		m.metrics.ErrorInc(errorType)
	}
}

type PublisherMetrics interface {
	MessageLag(duration float64, subject string)
	MessageCounterInc(subject string)
	ErrorInc(errorType string)
}

type publisherMetrics struct {
	metrics      ccmetrics.PublisherMetrics
	shouldReport bool
}

func NewPublisherMetrics(metrics ccmetrics.PublisherMetrics) PublisherMetrics {
	return publisherMetrics{
		metrics:      metrics,
		shouldReport: metrics != nil,
	}
}

func (m publisherMetrics) MessageLag(duration float64, subject string) {
	if m.shouldReport {
		m.metrics.MessageLag(duration, subject)
	}
}

func (m publisherMetrics) MessageCounterInc(subject string) {
	if m.shouldReport {
		m.metrics.MessageCounterInc(subject)
	}
}

func (m publisherMetrics) ErrorInc(errorType string) {
	if m.shouldReport {
		m.metrics.ErrorInc(errorType)
	}
}
